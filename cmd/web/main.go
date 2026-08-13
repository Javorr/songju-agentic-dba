package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/v1/devices", getAllDevices)
	mux.HandleFunc("GET /api/v1/devices/{id}", getDevice)
	mux.HandleFunc("POST /api/v1/devices/register", registerDevice)

	log.Print("Starting server on :8080")
	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}
