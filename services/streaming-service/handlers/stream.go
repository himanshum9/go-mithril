package handlers

import (
    "encoding/json"
    "net/http"
    "github.com/gorilla/mux"
    "log"
)

type StreamRequest struct {
    TenantID    string  `json:"tenant_id"`
    Latitude    float64 `json:"latitude"`
    Longitude   float64 `json:"longitude"`
}

func StreamLocationData(w http.ResponseWriter, r *http.Request) {
    var request StreamRequest

    err := json.NewDecoder(r.Body).Decode(&request)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    // Here you would typically process the location data and stream it to the third-party application.
    // For now, we'll just log the received data.
    log.Printf("Received location data: TenantID=%s, Latitude=%f, Longitude=%f", request.TenantID, request.Latitude, request.Longitude)

    // Respond with a success message
    w.WriteHeader(http.StatusAccepted)
    json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func RegisterRoutes(r *mux.Router) {
    r.HandleFunc("/stream", StreamLocationData).Methods("POST")
}