package middleware

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"fmt"

	"github.com/dgrijalva/jwt-go"
)

var (
	cognitoRegion     = os.Getenv("AWS_REGION")
	cognitoUserPoolID = os.Getenv("COGNITO_USER_POOL_ID")
)

func CognitoAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		token = strings.TrimPrefix(token, "Bearer ")
		claims, err := ValidateCognitoJWT(token)
		if err != nil {
			http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}
		ctx := r.Context()
		ctx = contextWithClaims(ctx, claims)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func ValidateCognitoJWT(tokenString string) (map[string]interface{}, error) {
	// Download JWKs from Cognito
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
	// Parse JWT and validate signature
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

// contextWithClaims and extraction helpers for RBAC
type contextKey string

func contextWithClaims(ctx context.Context, claims map[string]interface{}) context.Context {
	return context.WithValue(ctx, contextKey("claims"), claims)
}

func GetClaimsFromContext(ctx context.Context) map[string]interface{} {
	claims, _ := ctx.Value(contextKey("claims")).(map[string]interface{})
	return claims
}
