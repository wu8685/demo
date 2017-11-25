package main

import (
	"flag"
	"log"
	"os"

	"github.com/wu8685/demo/echo/server"
	"github.com/wu8685/demo/echo/server/config"

	_ "github.com/wu8685/demo/echo/handler/echo"
	_ "github.com/wu8685/demo/echo/handler/healthz"
	_ "github.com/wu8685/demo/echo/handler/tail"
)

func main() {
	f := flag.NewFlagSet("echo", flag.ExitOnError)
	config := config.New()
	f.IntVar(&config.Port, "port", config.Port, "port the server listens to")

	if err := f.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
	server.Run(config.Port)
}
