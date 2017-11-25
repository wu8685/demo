package tools

import (
	"encoding/json"
	"github.com/cloudflare/cfssl/log"
	"net/http"
)

func WriteResponse(w http.ResponseWriter, statusCode int, body interface{}, err error) {
	if err != nil {
		WriteError(w, err)
		return
	}

	bs, err := json.Marshal(body)
	if err != nil {
		WriteError(w, err)
		return
	}

	w.WriteHeader(statusCode)
	w.Write(bs)
}

func WriteError(w http.ResponseWriter, err error) {
	log.Error("response err: %s", err)
	w.WriteHeader(400)
	w.Write([]byte(err.Error()))
}
