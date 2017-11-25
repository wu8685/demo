package main

import (
	"log"
	"os"
	"time"
)

func main() {
	timer := 5 * time.Second
	log.Printf("Will exit(1) after %v.\n", timer)

	<-time.After(timer)
	log.Println("Exit(1) as expected")
	os.Exit(1)
}
