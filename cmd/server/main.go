package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	metrics "github.com/Pablo-Loyola-Tantaruna/go-edge-gateway/internal"
	"github.com/Pablo-Loyola-Tantaruna/go-edge-gateway/internal/balancer"
	"github.com/Pablo-Loyola-Tantaruna/go-edge-gateway/internal/config"
	"github.com/Pablo-Loyola-Tantaruna/go-edge-gateway/internal/health"
	"github.com/Pablo-Loyola-Tantaruna/go-edge-gateway/internal/proxy"
	"github.com/Pablo-Loyola-Tantaruna/go-edge-gateway/internal/proxy/middleware"
	"github.com/Pablo-Loyola-Tantaruna/go-edge-gateway/pkg/models"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	b := config.InitBootstrap()
	log.Printf("Iniciando %s en ambiente: %s", b.AppName, b.Environment)

	cfg, err := config.LoadConfig(b.ConfigPath)

	jwtSecret := os.Getenv("GOPHER_JWT_SECRET")
	if jwtSecret == "" {
		if b.Environment == "prod" {
			log.Fatal("CRITICAL: JWT_SECRET no encontrado en PROD")
		}
		jwtSecret = "dev_secret"
	}

	if err != nil {
		log.Fatalf("Error cargando configuración: %v", err)
	}

	pool := &balancer.ServerPool{}

	for _, b := range cfg.Backends {
		serverURL, _ := url.Parse(b.URL)
		proxyBackend := httputil.NewSingleHostReverseProxy(serverURL)

		pool.AddBackend(&models.Backend{
			URL:          serverURL,
			Alive:        true,
			ReverseProxy: proxyBackend,
		})

		proxyBackend.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			log.Printf("Error conectando al backend %s: %v. Reintentando...", serverURL.Host, err)

			retries := proxy.GetRetryFromContext(r)

			if retries < 3 {

				log.Printf("Fallo en backend %s. Reintento %d/3", serverURL.Host, retries+1)
				nextPeer := pool.GetNextPeerAfterFailure()

				if nextPeer != nil {
					r = proxy.SetRetryInContext(r, retries+1)
					nextPeer.ReverseProxy.ServeHTTP(w, r)
					return
				}
			}
			http.Error(w, "Servicio no disponible (No backends alive)", http.StatusServiceUnavailable)
		}
	}
	config.WatchConfig("config.yaml", func(newCfg *config.Config) {
		var newBackends []*models.Backend
		for _, b := range newCfg.Backends {
			serverURL, _ := url.Parse(b.URL)
			proxyBackend := httputil.NewSingleHostReverseProxy(serverURL)

			newBackends = append(newBackends, &models.Backend{
				URL:          serverURL,
				Alive:        true,
				ReverseProxy: proxyBackend,
			})
		}
		pool.UpdateBackends(newBackends)
		log.Printf(" Configuración actualizada: %d servidores activos", len(newBackends))
	})

	go health.HealthCheck(pool)
	http.Handle("/metrics", promhttp.Handler())

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	originHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		r = r.WithContext(ctx)

		peer := pool.GetNextPeer()

		if peer != nil {
			peer.ReverseProxy.ServeHTTP(w, r)
			return
		}

		http.Error(w, "Servicio no disponible (No backends alive)", http.StatusServiceUnavailable)
	})

	limiter := middleware.NewIPRateLimiter(5, 10)

	authHandler := middleware.JWTMiddleware(jwtSecret, originHandler)
	limiterHandler := middleware.RateLimitMiddleware(limiter, authHandler)
	finalHandler := Logger(limiterHandler)

	server := http.Server{
		Addr:    addr,
		Handler: finalHandler,
	}

	log.Printf("GopherGuard iniciado en el puerto %s", addr)
	log.Fatal(server.ListenAndServe())
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = fmt.Sprintf("%d", time.Now().UnixNano())
			r.Header.Set("X-Request-ID", reqID)
		}

		w.Header().Set("X-Request-ID", reqID)

		next.ServeHTTP(w, r)

		duration := time.Since(start).Seconds()
		metrics.HttpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, "200").Inc()
		metrics.RequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)

		log.Printf("[%s] %s %s | Duración: %v", reqID, r.Method, r.URL.Path, time.Since(start))
	})
}
