package echo

import (
	"bytes"
	"log"
	"net/http"

	"github.com/wu8685/demo/echo/server"
	"github.com/wu8685/demo/echo/tools"
)

func init() {
	server.Register(echo, "/echo", "POST")
}

func echo(w http.ResponseWriter, r *http.Request) {
	log.Printf("start handling echo request")

	body := &bytes.Buffer{}
	_, err := body.ReadFrom(r.Body)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
	}
	log.Printf("handle echo: %s", body.String())
	tools.WriteResponse(w, 200, body, nil)
}
