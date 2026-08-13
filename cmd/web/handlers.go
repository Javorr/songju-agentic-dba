package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Javorr/telemetry-platform/internal/dto"
)

var (
	devices      = map[string]*dto.Device{}
	nextDeviceID = 1
)

func getAllDevices(w http.ResponseWriter, r *http.Request) {
	list := make([]*dto.Device, 0, len(devices))
	for _, d := range devices {
		list = append(list, d)
	}

	writeJSON(w, http.StatusOK, list)
}

func getDevice(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	d, ok := devices[id]
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "device not found"})
		return
	}

	writeJSON(w, http.StatusOK, d)
}

func registerDevice(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if req.DeviceID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "deviceId is required"})
		return
	}

	if _, exists := devices[req.DeviceID]; exists {
		writeJSON(w, http.StatusConflict, map[string]string{"error": "device already registered"})
		return
	}

	token := generateToken()

	d := &dto.Device{
		ID:              nextDeviceID,
		DeviceID:        req.DeviceID,
		FirmwareVersion: req.FirmwareVersion,
		RegisteredAt:    time.Now(),
		LastSeen:        time.Now(),
		Status:          "online",
		ReportInterval:  60,
		Token:           token,
	}
	nextDeviceID++
	devices[req.DeviceID] = d

	log.Printf("device registered: %s", req.DeviceID)

	writeJSON(w, http.StatusCreated, d)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func generateToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
