package models

type Tenant struct {
    TenantID   string `json:"tenant_id"`
    Name       string `json:"name"`
    Description string `json:"description"`
}

func NewTenant(tenantID, name, description string) *Tenant {
    return &Tenant{
        TenantID:   tenantID,
        Name:       name,
        Description: description,
    }
}