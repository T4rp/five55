package main

import (
	"log"
	"os"
)

func main() {
	token := os.Getenv("TOKEN")

	if token == "" {
		log.Fatal("No token provided")
	}
}
