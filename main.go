package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"syscall"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// LogPayload represents the incoming data from DB/Apps
type LogPayload struct {
	SQL      string `json:"sql"`
	Service  string `json:"service"`
	Duration int    `json:"duration_ms"`
}

func main() {
	// Connect to NATS and create JetStream for allowing persistence later
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	js, _ := jetstream.New(nc)

	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()
	// Create a Stream named "DB_LOGS" to persist our messages
	_, err = js.CreateStream(ctx, jetstream.StreamConfig{
		Name:     "DB_LOGS",
		Subjects: []string{"db.logs.slow"},
	})
	if err != nil {
		log.Printf("Stream might already exist: %v", err)
	}

	// HTTP Handler
	http.HandleFunc("/ingest", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST is allowed", http.StatusMethodNotAllowed)
			return
		}

		var payload LogPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Mask emails and UUIDs before they ever touch the message bus
		payload.SQL = maskPII(payload.SQL)

		// Publish to NATS
		data, _ := json.Marshal(payload)
		_, err := js.Publish(r.Context(), "db.logs.slow", data)
		if err != nil {
			log.Printf("Failed to publish: %v", err)
			http.Error(w, "Queue failure", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Log queued for analysis")
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}
	log.Println("Ingestor running on " + port)
	go func() {
		if err := http.ListenAndServe(port, nil); err != nil && err != http.ErrServerClosed {
			log.Printf("Server error: %v", err)
		}
	}()
}

var (
	emailRegex      = regexp.MustCompile(`'[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,4}'`)
	idRegex         = regexp.MustCompile(`(=|IN)\s*\(?'?[0-9]+'?(,\s*'?[0-9]+'?)*\)?`)
	ssnRegex        = regexp.MustCompile(`'\d{3}-\d{2}-\d{4}'`)
	creditCardRegex = regexp.MustCompile(`'\d{4}-\d{4}-\d{4}-\d{4}'`)
	uuidRegex       = regexp.MustCompile(`(?i)'[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}'`)
)

// maskPII is a fast utility to scrub sensitive data from SQL strings
func maskPII(sql string) string {
	sql = emailRegex.ReplaceAllString(sql, "[REDACTED_EMAIL]")
	sql = ssnRegex.ReplaceAllString(sql, "[REDACTED_SSN]")
	sql = creditCardRegex.ReplaceAllString(sql, "[REDACTED_CARD]")
	sql = uuidRegex.ReplaceAllString(sql, "[REDACTED_UUID]")
	sql = idRegex.ReplaceAllString(sql, "$1 [REDACTED_ID]")

	return sql
}
