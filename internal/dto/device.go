package dto

import "time"

type Device struct {
	ID              int       `json:"id"`
	DeviceID        string    `json:"deviceId"`
	FirmwareVersion string    `json:"firmwareVersion"`
	RegisteredAt    time.Time `json:"registeredAt"`
	LastSeen        time.Time `json:"lastSeen"`
	Status          string    `json:"status"`
	ReportInterval  int       `json:"reportInterval"`
	Token           string    `json:"token"`
}

type RegisterRequest struct {
	DeviceID        string `json:"deviceId"`
	FirmwareVersion string `json:"firmwareVersion"`
}
