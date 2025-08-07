package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/himanshum9/go-mithril/services/location-service/models"
)

var (
	cognitoRegion     = os.Getenv("AWS_REGION")
	cognitoUserPoolID = os.Getenv("COGNITO_USER_POOL_ID")
)

func main() {
	connStr := os.Getenv("LOCATION_DB_CONN")
	if connStr == "" {
		connStr = "host=localhost port=5432 user=your_db_user password=your_db_password dbname=multi_tenant_db sslmode=disable"
	}
	if err := models.InitDB(connStr); err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}

	router := gin.Default()
	router.Use(CognitoAuthMiddleware)
	router.POST("/location", submitLocation)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

func CognitoAuthMiddleware(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	token = strings.TrimPrefix(token, "Bearer ")
	claims, err := ValidateCognitoJWT(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: " + err.Error()})
		return
	}
	c.Set("claims", claims)
	c.Next()
}

func ValidateCognitoJWT(tokenString string) (map[string]interface{}, error) {
	jwksURL := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json", cognitoRegion, cognitoUserPoolID)
	resp, err := http.Get(jwksURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var jwks struct {
		Keys []map[string]interface{} `json:"keys"`
	}
	if err := json.Unmarshal(body, &jwks); err != nil {
		return nil, err
	}
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("no kid in token header")
		}
		for _, key := range jwks.Keys {
			if key["kid"] == kid {
				// TODO: Parse JWK to PEM and return *rsa.PublicKey
				return nil, fmt.Errorf("JWK to PEM conversion not implemented")
			}
		}
		return nil, fmt.Errorf("public key not found")
	})
	if err != nil || !parsedToken.Valid {
		return nil, fmt.Errorf("invalid token: %v", err)
	}
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}
	return claims, nil
}

func submitLocation(c *gin.Context) {
	var l models.Location
	if err := c.ShouldBindJSON(&l); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := models.SaveLocation(context.Background(), &l); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Location submitted successfully", "location": l})
}
