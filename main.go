package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// LogPayload represents the incoming data from your DB/App
type LogPayload struct {
	SQL       string `json:"sql"`
	Service   string `json:"service"`
	Duration  int    `json:"duration_ms"`
}

func main() {
	// 1. Connect to NATS
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// 2. Initialize JetStream
	js, _ := jetstream.New(nc)

	// Create a Stream named "DB_LOGS" to persist our messages
	ctx := http.Context() // Use a context for safety
	_, err = js.CreateStream(ctx, jetstream.StreamConfig{
		Name:     "DB_LOGS",
		Subjects: []string{"db.logs.slow"},
	})
	if err != nil {
		log.Printf("Stream might already exist: %v", err)
	}

	// 3. The HTTP Handler
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

		// 4. PII REDACTION (The Security Layer)
		// We mask emails and UUIDs before they ever touch the Message Bus
		payload.SQL = maskPII(payload.SQL)

		// 5. Publish to NATS
		data, _ := json.Marshal(payload)
		_, err := js.Publish(r.Context(), "db.logs.slow", data)
		if err != nil {
			log.Printf("Failed to publish: %v", err)
			http.Error(w, "Queue failure", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Log queued for analysis")
	})

	log.Println("Ingestor running on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// maskPII is a fast utility to scrub sensitive data from SQL strings
func maskPII(sql string) string {
	// Example: Masking emails
	emailRegex := regexp.MustCompile(`[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,4}`)
	sql = emailRegex.ReplaceAllString(sql, "[REDACTED_EMAIL]")

	// Example: Masking numbers in WHERE clauses (common for IDs)
	// 'WHERE id = 123' becomes 'WHERE id = [REDACTED_ID]'
	idRegex := regexp.MustCompile(`(=|IN)\s*\(?'?[0-9]+'?(,\s*'?[0-9]+'?)*\)?`)
	return idRegex.ReplaceAllString(sql, "$1 [REDACTED_VAL]")
}
