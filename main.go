package main

import (
	"log"
	"os"
)

func main() {
	server, err := NewServer(
		os.Getenv("C_BACKEND_URL"),
		os.Getenv("A_BACKEND_URL"),
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(server.StartApp())
}
