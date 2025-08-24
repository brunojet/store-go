package main

import (
	"log"

	"github.com/brunojet/store-go/office/api"
)

func main() {
	if err := api.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
