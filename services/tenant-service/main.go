package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/himanshum9/go-mithril/services/tenant-service/models"
)

func main() {
	connStr := os.Getenv("TENANT_DB_CONN")
	if connStr == "" {
		connStr = "host=localhost port=5432 user=your_db_user password=your_db_password dbname=multi_tenant_db sslmode=disable"
	}
	if err := models.InitDB(connStr); err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/tenants", createTenant).Methods("POST")
	router.HandleFunc("/tenants/{id}", getTenant).Methods("GET")
	router.HandleFunc("/tenants", listTenants).Methods("GET")

	log.Println("Starting tenant service on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}

// Placeholder functions for tenant management
func createTenant(w http.ResponseWriter, r *http.Request) {
	var t models.Tenant
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := models.SaveTenant(context.Background(), &t); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(t)
}

func getTenant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenantID := vars["id"]
	t, err := models.GetTenant(context.Background(), tenantID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(t)
}

func listTenants(w http.ResponseWriter, r *http.Request) {
	tenants, err := models.ListTenants(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(tenants)
}
