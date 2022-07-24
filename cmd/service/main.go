package main

import (
	"log"

	"github.com/ineverbee/wbl0/internal/app"
)

func main() {
	err := app.StartApp()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Server stopped!")
}
