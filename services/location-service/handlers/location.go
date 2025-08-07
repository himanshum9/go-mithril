package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/himanshum9/go-mithril/services/location-service/models"
	"github.com/himanshum9/go-mithril/services/location-service/streaming"
)

// LocationSubmissionRequest represents the payload for location submission
type LocationSubmissionRequest struct {
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
	SessionID string  `json:"session_id" binding:"required"`
	Timestamp int64   `json:"timestamp" binding:"required"`
}

// In-memory session tracker (for demo; use Redis/DB in prod)
var sessionLastSubmission = make(map[string]int64)

var (
	kafkaBroker   string
	kafkaStreamer *streaming.Streamer
)

func init() {
	kafkaBroker = os.Getenv("KAFKA_BROKER")
	if kafkaBroker == "" {
		kafkaBroker = "localhost:9092"
	}
	kafkaStreamer = streaming.NewStreamer(kafkaBroker, "location-stream")
}

// SubmitLocation handles location data submission by tenant users
func SubmitLocation(w http.ResponseWriter, r *http.Request) {
	// Extract JWT claims from context (pseudo-code, adapt to your middleware)
	claims, ok := r.Context().Value("claims").(map[string]interface{})
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	role, _ := claims["role"].(string)
	tenantID, _ := claims["tenant_id"].(string)
	if role != "user" {
		http.Error(w, "Forbidden: Tenant users only", http.StatusForbidden)
		return
	}

	var req LocationSubmissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Session timing: allow only 10 min sessions
	startTime := req.Timestamp - (req.Timestamp % 600)
	if req.Timestamp < startTime || req.Timestamp > startTime+600 {
		http.Error(w, "Session expired or invalid timestamp", http.StatusBadRequest)
		return
	}

	// Interval check: allow only every 30 seconds
	last, ok := sessionLastSubmission[req.SessionID]
	if ok && req.Timestamp-last < 30 {
		http.Error(w, "Submission interval too short", http.StatusTooManyRequests)
		return
	}
	sessionLastSubmission[req.SessionID] = req.Timestamp

	timestamp := time.Unix(req.Timestamp, 0)
	location := models.Location{
		TenantID:  tenantID,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		Timestamp: timestamp,
	}

	// Save location to DB
	if err := models.SaveLocation(r.Context(), &location); err != nil {
		http.Error(w, "Failed to save location", http.StatusInternalServerError)
		return
	}

	// Stream location to Kafka
	err := kafkaStreamer.StreamLocationData(location.TenantID, location.Latitude, location.Longitude, location.Timestamp)
	if err != nil {
		// Attempt to re-establish connection and retry once
		kafkaStreamer.Close()
		kafkaStreamer = streaming.NewStreamer(kafkaBroker, "location-stream")
		err = kafkaStreamer.StreamLocationData(location.TenantID, location.Latitude, location.Longitude, location.Timestamp)
		if err != nil {
			http.Error(w, "Failed to stream location data", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "Location submitted and streamed", "location": location})
}

// ...existing code...
