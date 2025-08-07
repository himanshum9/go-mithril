package handlers

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/himanshum9/go-mithril/services/tenant-service/models"
)

// CreateTenant handles the creation of a new tenant
func CreateTenant(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userClaims := claims.(*jwt.MapClaims)
	role := (*userClaims)["role"].(string)
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: Admins only"})
		return
	}

	var tenant models.Tenant
	if err := c.ShouldBindJSON(&tenant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save tenant to DB, ensure tenant_id is unique
	// ...existing code...

	c.JSON(http.StatusCreated, tenant)
}

// GetTenant retrieves a tenant by ID
func GetTenant(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userClaims := claims.(*jwt.MapClaims)
	role := (*userClaims)["role"].(string)
	tenantID := c.Param("id")
	userTenantID := (*userClaims)["tenant_id"].(string)
	if role != "admin" && tenantID != userTenantID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: Access denied"})
		return
	}

	// Retrieve tenant from DB by tenantID
	var tenant models.Tenant // Assume we fetched the tenant from the database
	c.JSON(http.StatusOK, tenant)
}

// ListTenants retrieves all tenants
func ListTenants(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userClaims := claims.(*jwt.MapClaims)
	role := (*userClaims)["role"].(string)
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: Admins only"})
		return
	}

	// Retrieve all tenants from DB
	var tenants []models.Tenant // Assume we fetched the tenants from the database
	c.JSON(http.StatusOK, tenants)
}
