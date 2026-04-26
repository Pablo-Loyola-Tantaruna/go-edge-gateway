package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	metrics "github.com/Pablo-Loyola-Tantaruna/go-edge-gateway/internal"
	"github.com/Pablo-Loyola-Tantaruna/go-edge-gateway/internal/balancer"
	"github.com/Pablo-Loyola-Tantaruna/go-edge-gateway/internal/config"
	"github.com/Pablo-Loyola-Tantaruna/go-edge-gateway/internal/health"
	"github.com/Pablo-Loyola-Tantaruna/go-edge-gateway/pkg/models"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Error cargando configuración: %v", err)
	}

	pool := &balancer.ServerPool{}

	for _, b := range cfg.Backends {
		serverURL, _ := url.Parse(b.URL)
		proxy := httputil.NewSingleHostReverseProxy(serverURL)

		pool.AddBackend(&models.Backend{
			URL:          serverURL,
			Alive:        true,
			ReverseProxy: proxy,
		})
	}

	go health.HealthCheck(pool)
	http.Handle("/metrics", promhttp.Handler())

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	originHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		peer := pool.GetNextPeer()

		if peer != nil {
			peer.ReverseProxy.ServeHTTP(w, r)
			return
		}

		http.Error(w, "Servicio no disponible (No backends alive)", http.StatusServiceUnavailable)
	})
	
	server := http.Server{
		Addr:    addr,
		Handler: Logger(originHandler),
	}

	log.Printf("🚀 GopherGuard iniciado en el puerto %s", addr)
	log.Fatal(server.ListenAndServe())
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		duration := time.Since(start).Seconds()

		metrics.HttpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, "200").Inc()
		metrics.RequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)

		log.Printf("Método: %s | Ruta: %s | Duración: %v", r.Method, r.URL.Path, time.Since(start))
	})
}
