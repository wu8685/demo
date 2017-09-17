package healthz

import (
	"net/http"

	"demo/echo/server"
	"log"
)

func init() {
	server.Register(healthz, "/healthz", "GET")
}

func healthz(w http.ResponseWriter, r *http.Request) {
	log.Printf("handle health check")
	w.Write([]byte("OK"))
}
