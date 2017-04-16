package healthz

import (
	"net/http"

	"demo/echo/server"
)

func init() {
	server.Register(healthz, "/healthz", "GET")
}

func healthz(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
