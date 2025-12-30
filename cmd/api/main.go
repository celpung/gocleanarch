package main

import (
	"log"
	"net/http"

	"github.com/celpung/gocleanarch/internal/infra/container"
)

func main() {
	app, err := container.Build()
	if err != nil {
		log.Fatalf("failed to boot application: %v", err)
	}

	port := app.Env.PORT
	log.Printf("Server running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, app.Router))
}
