package kontaktauth

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
)

type JWTAuth struct {
	KeycloakURL   string
	jwksCache     map[string]*keyfunc.JWKS
	jwksCacheLock sync.RWMutex
	jwksOpts      keyfunc.Options
}

func NewJWTAuth() *JWTAuth {
	opts := keyfunc.Options{
		RefreshErrorHandler: func(err error) {
			log.Printf("JWKS refresh error: %v", err)
		},
		RefreshInterval:   time.Hour,
		RefreshUnknownKID: true,
	}
	return &JWTAuth{
		jwksCache: make(map[string]*keyfunc.JWKS),
		jwksOpts:  opts,
	}
}

func (ja *JWTAuth) ExtractCompanyID(tokenStr string) (string, error) {
	if strings.HasPrefix(strings.ToLower(tokenStr), "bearer ") {
		tokenStr = tokenStr[len("Bearer "):]
	}

	// Parse unverified to get issuer
	parser := new(jwt.Parser)
	unverified, _, err := parser.ParseUnverified(tokenStr, jwt.MapClaims{})
	if err != nil {
		return "", fmt.Errorf("cannot parse token header: %w", err)
	}
	claimsU := unverified.Claims.(jwt.MapClaims)
	iss, ok := claimsU["iss"].(string)
	if !ok {
		return "", errors.New("issuer claim missing or not a string")
	}

	// Validate issuer prefix
	prefix := strings.TrimRight(ja.KeycloakURL, "/") + "/"
	if !strings.HasPrefix(iss, prefix) {
		return "", fmt.Errorf("invalid issuer %q: must start with %q", iss, prefix)
	}
	realmPath := strings.TrimPrefix(iss, prefix)
	realm := strings.SplitN(realmPath, "/", 2)[0]

	// Load or fetch JWKS for realm
	ja.jwksCacheLock.RLock()
	jwks := ja.jwksCache[realm]
	ja.jwksCacheLock.RUnlock()
	if jwks == nil {
		url := fmt.Sprintf("%s/%s/protocol/openid-connect/certs", ja.KeycloakURL, realm)
		jwksNew, err := keyfunc.Get(url, ja.jwksOpts)
		if err != nil {
			return "", fmt.Errorf("failed to load JWKS for realm %q: %w", realm, err)
		}
		ja.jwksCacheLock.Lock()
		ja.jwksCache[realm] = jwksNew
		ja.jwksCacheLock.Unlock()
		jwks = jwksNew
	}

	// Verify signature, claims exp/nbf
	token, err := jwt.ParseWithClaims(tokenStr, jwt.MapClaims{}, jwks.Keyfunc)
	if err != nil {
		return "", fmt.Errorf("jwt parse/verify error: %w", err)
	}
	if !token.Valid {
		return "", errors.New("token expired or not valid yet")
	}

	// Extract company-id claim
	claims := token.Claims.(jwt.MapClaims)
	rawCid, ok := claims["company-id"]
	if !ok {
		return "", errors.New("company-id claim missing")
	}
	cid, ok := rawCid.(string)
	if !ok {
		return "", errors.New("company-id claim is not a string")
	}
	return cid, nil
}
