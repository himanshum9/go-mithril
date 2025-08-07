package models

import (
	"context"
	"time"
)

func SaveLocation(ctx context.Context, l *Location) error {
	_, err := DB.ExecContext(ctx, `INSERT INTO locations (latitude, longitude, timestamp, tenant_id) VALUES ($1, $2, $3, $4)`, l.Latitude, l.Longitude, l.Timestamp, l.TenantID)
	return err
}

func ListLocationsByTenant(ctx context.Context, tenantID string) ([]Location, error) {
	rows, err := DB.QueryContext(ctx, `SELECT latitude, longitude, timestamp, tenant_id FROM locations WHERE tenant_id = $1`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var locations []Location
	for rows.Next() {
		var l Location
		var ts time.Time
		if err := rows.Scan(&l.Latitude, &l.Longitude, &ts, &l.TenantID); err != nil {
			return nil, err
		}
		l.Timestamp = ts
		locations = append(locations, l)
	}
	return locations, nil
}
