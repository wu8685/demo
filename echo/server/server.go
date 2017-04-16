package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type RestInfo struct {
	Path    string
	Methods []string
	Handler func(http.ResponseWriter, *http.Request)
}

func Register(handler func(http.ResponseWriter, *http.Request), path string, methods ...string) {
	rests = append(rests, &RestInfo{path, methods, handler})
}

var rests = []*RestInfo{}

func Run(port int) {
	router := mux.NewRouter()
	for _, info := range rests {
		log.Printf("Register REST API: %s - %v\n", info.Path, info.Methods)
		router.Path(info.Path).Methods(info.Methods...).HandlerFunc(info.Handler)
	}

	addr := fmt.Sprintf(":%d", port)
	log.Printf("Start on %s\n", addr)
	http.ListenAndServe(addr, router)
}
