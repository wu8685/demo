package main

import (
	"flag"
	"os"
	"log"

	"demo/echo/server"
	"demo/echo/server/config"

	_ "demo/echo/handler/echo"
	_ "demo/echo/handler/healthz"
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
