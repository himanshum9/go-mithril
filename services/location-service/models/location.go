package models

import "time"

// Location represents the geographical location data submitted by a tenant user.
type Location struct {
    Latitude  float64   `json:"latitude"`
    Longitude float64   `json:"longitude"`
    Timestamp time.Time `json:"timestamp"`
    TenantID  string    `json:"tenant_id"` // Identifier for the tenant to which this location data belongs
}