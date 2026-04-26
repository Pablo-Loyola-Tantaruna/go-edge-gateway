package health

import (
	"github.com/Pablo-Loyola-Tantaruna/go-edge-gateway/internal/balancer"
	_ "github.com/Pablo-Loyola-Tantaruna/go-edge-gateway/pkg/models"
	"log"
	"net"
	"net/url"
	"time"
)

func HealthCheck(s *balancer.ServerPool) {
	t := time.NewTicker(time.Second * 20) // Cada 20 segundos
	for {
		select {
		case <-t.C:
			log.Println("Iniciando verificación de salud...")
			for _, b := range s.Backends {
				alive := isBackendAlive(b.URL)
				b.SetAlive(alive)
				status := "activo"
				if !alive {
					status = "caído"
				}
				log.Printf("Servidor %s está %s\n", b.URL, status)
			}
		}
	}
}

func isBackendAlive(u *url.URL) bool {
	timeout := 2 * time.Second
	conn, err := net.DialTimeout("tcp", u.Host, timeout)
	if err != nil {
		log.Printf("Error de salud en %s: %v\n", u.Host, err)
		return false
	}
	err = conn.Close()
	if err != nil {
		return false
	}
	return true
}
