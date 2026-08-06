package main

import (
	"log"
	"net/http"
)

func getAllDevices(w http.ResponseWriter, r *http.Request) {
	log.Print("GET devices")
}

func getDevice(w http.ResponseWriter, r *http.Request) {
	log.Print("GET device")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/devices", getAllDevices)
	mux.HandleFunc("/api/v1/device/", getDevice)

	log.Print("Starting server on :8080")
	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}
