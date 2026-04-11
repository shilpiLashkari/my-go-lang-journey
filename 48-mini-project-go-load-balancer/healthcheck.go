package main

import (
	"log"
	"net"
	"net/url"
	"time"
)

// isBackendAlive checks whether a backend is Alive by establishing a TCP connection
func isBackendAlive(u *url.URL) bool {
	timeout := 2 * time.Second
	conn, err := net.DialTimeout("tcp", u.Host, timeout)
	if err != nil {
		log.Println("Site unreachable, error: ", err)
		return false
	}
	defer conn.Close()
	return true
}

// healthCheck periodically checks the health of the backends
func healthCheck(s *ServerPool) {
	t := time.NewTicker(time.Second * 30)
	for {
		select {
		case <-t.C:
			log.Println("Starting health check...")
			s.HealthCheck()
			log.Println("Health check completed")
		}
	}
}
