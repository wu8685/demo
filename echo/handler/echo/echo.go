package echo

import (
	"bytes"
	"net/http"

	"demo/echo/server"
)

func init() {
	server.Register(echo, "/echo", "POST")
}

func echo(w http.ResponseWriter, r *http.Request) {
	body := &bytes.Buffer{}
	_, err := body.ReadFrom(r.Body)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
	}
	w.Write(body.Bytes())
}
