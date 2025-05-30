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

type TokenValidator interface {
	ValidateToken(tokenStr string) bool
}

type JwksValidator struct {
	KeycloakURL   string
	Audience      string
	jwksCache     map[string]*keyfunc.JWKS
	jwksCacheLock sync.RWMutex
	jwksOpts      keyfunc.Options
}

func (ja *JwksValidator) ValidateToken(tokenStr string) bool {
	companyId, _ := ExtractCompanyID(tokenStr)
	log.Printf("[jwksValidator] validating token for company %s", companyId)
	realm, err := ja.verifyIss(tokenStr)
	if err != nil {
		log.Printf("[jwksValidator] verifyIss error: %v", err)
		return false
	}

	jwks, err := ja.getJWKS(realm)
	if err != nil {
		log.Printf("[jwksValidator] getJWKS error for realm %s: %v", realm, err)
		return false
	}

	token, err := jwt.ParseWithClaims(tokenStr, jwt.MapClaims{}, jwks.Keyfunc)
	if err != nil {
		log.Printf("[jwksValidator] jwt.Parse error: %v", err)
		return false
	}
	if !token.Valid {
		log.Printf("[jwksValidator] token not valid")
		return false
	}

	claims := token.Claims.(jwt.MapClaims)
	if err := ja.verifyAud(claims); err != nil {
		log.Printf("[jwksValidator] verifyAud error: %v", err)
		return false
	}

	return true
}

func (ja *JwksValidator) verifyIss(tokenStr string) (string, error) {
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

	prefix := ja.KeycloakURL
	if !strings.HasPrefix(iss, prefix) {
		return "", fmt.Errorf("invalid issuer %q: must start with %q", iss, prefix)
	}
	realm := strings.SplitN(strings.TrimPrefix(iss, prefix), "/", 2)[0]
	return realm, nil
}

func (ja *JwksValidator) getJWKS(realm string) (*keyfunc.JWKS, error) {
	ja.jwksCacheLock.RLock()
	jwks := ja.jwksCache[realm]
	ja.jwksCacheLock.RUnlock()
	if jwks != nil {
		return jwks, nil
	}

	url := fmt.Sprintf("%s%s/protocol/openid-connect/certs", ja.KeycloakURL, realm)
	log.Printf("[JwksValidator] fetching JWKS from URL=%s", url)
	jwksNew, err := keyfunc.Get(url, ja.jwksOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to load JWKS for realm %q: %w", realm, err)
	}
	ja.jwksCacheLock.Lock()
	ja.jwksCache[realm] = jwksNew
	ja.jwksCacheLock.Unlock()
	return jwksNew, nil
}

func (ja *JwksValidator) verifyAud(claims jwt.MapClaims) error {
	if ja.Audience == "" {
		return nil
	}
	audVal, ok := claims["aud"]
	if !ok {
		return errors.New("aud claim missing")
	}
	switch v := audVal.(type) {
	case string:
		if v != ja.Audience {
			return fmt.Errorf("invalid audience %q", v)
		}
	case []interface{}:
		for _, vv := range v {
			if s, ok := vv.(string); ok && s == ja.Audience {
				return nil
			}
		}
		return fmt.Errorf("audience %q not present in token", ja.Audience)
	default:
		return errors.New("aud claim has unknown type")
	}
	return nil
}

type cacheEntry struct {
	valid bool
	exp   time.Time
}

type CachingValidator struct {
	base      TokenValidator
	mu        sync.RWMutex
	cache     map[string]*cacheEntry
	jwtParser *jwt.Parser
}

func (c *CachingValidator) ValidateToken(tokenStr string) bool {
	key := extractSignature(tokenStr)

	c.mu.RLock()
	entry, ok := c.cache[key]
	c.mu.RUnlock()
	if ok {
		if entry.valid {
			if time.Now().After(entry.exp) {
				entry.valid = false
				log.Printf("[cachingValidator] token expired for signature %s", key)
				return false
			}
			return true
		}
		return false
	}

	valid := c.base.ValidateToken(tokenStr)
	expTime := time.Now()

	tokUnv, _, err := c.jwtParser.ParseUnverified(tokenStr, jwt.MapClaims{})
	if err == nil {
		if cl, ok2 := tokUnv.Claims.(jwt.MapClaims); ok2 {
			if expVal, ok3 := cl["exp"].(float64); ok3 {
				expTime = time.Unix(int64(expVal), 0)
			}
		}
	}

	c.mu.Lock()
	c.cache[key] = &cacheEntry{valid: valid, exp: expTime}
	c.mu.Unlock()

	return valid
}

func extractSignature(tokenStr string) string {
	parts := strings.Split(tokenStr, ".")
	if len(parts) == 3 {
		return parts[2]
	}
	return tokenStr
}
