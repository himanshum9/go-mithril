package models

import (
	"context"
)

func SaveTenant(ctx context.Context, t *Tenant) error {
	_, err := DB.ExecContext(ctx, `INSERT INTO tenants (tenant_id, name, description) VALUES ($1, $2, $3)`, t.TenantID, t.Name, t.Description)
	return err
}

func GetTenant(ctx context.Context, tenantID string) (*Tenant, error) {
	row := DB.QueryRowContext(ctx, `SELECT tenant_id, name, description FROM tenants WHERE tenant_id = $1`, tenantID)
	var t Tenant
	if err := row.Scan(&t.TenantID, &t.Name, &t.Description); err != nil {
		return nil, err
	}
	return &t, nil
}

func ListTenants(ctx context.Context) ([]Tenant, error) {
	rows, err := DB.QueryContext(ctx, `SELECT tenant_id, name, description FROM tenants`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tenants []Tenant
	for rows.Next() {
		var t Tenant
		if err := rows.Scan(&t.TenantID, &t.Name, &t.Description); err != nil {
			return nil, err
		}
		tenants = append(tenants, t)
	}
	return tenants, nil
}
