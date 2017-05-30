package server

import (
	"fmt"
	"log"
	"net/http"
	"net/http/pprof"

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
	attachProfiler(router)

	addr := fmt.Sprintf(":%d", port)
	log.Printf("Start on %s\n", addr)
	http.ListenAndServe(addr, router)
}

func attachProfiler(router *mux.Router) {
	router.HandleFunc("/debug/pprof/", pprof.Index)
	router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)

	// Manually add support for paths linked to by index page at /debug/pprof/
	router.Handle("/debug/pprof/block", pprof.Handler("block"))
	router.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	router.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	router.Handle("/debug/pprof/mutex", pprof.Handler("mutex"))
	router.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
}
