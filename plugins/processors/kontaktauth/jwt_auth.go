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
	Audience      string
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
		Audience:  "compute-api",
		jwksCache: make(map[string]*keyfunc.JWKS),
		jwksOpts:  opts,
	}
}

func (ja *JWTAuth) ExtractCompanyID(claims jwt.MapClaims) (string, error) {
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

func (ja *JWTAuth) VerifyToken(tokenStr string) (jwt.MapClaims, error) {
	if strings.HasPrefix(strings.ToLower(tokenStr), "bearer ") {
		tokenStr = tokenStr[len("Bearer "):]
	}

	realm, err := ja.verifyIss(tokenStr)
	if err != nil {
		return nil, err
	}

	jwks, err := ja.getJWKS(realm)
	if err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(tokenStr, jwt.MapClaims{}, jwks.Keyfunc)
	if err != nil {
		return nil, fmt.Errorf("jwt parse/verify error: %w", err)
	}
	if !token.Valid {
		return nil, errors.New("token expired or not valid yet")
	}

	claims := token.Claims.(jwt.MapClaims)

	if err := ja.verifyAud(claims); err != nil {
		return nil, err
	}

	return claims, nil
}

func (ja *JWTAuth) verifyIss(tokenStr string) (string, error) {
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

func (ja *JWTAuth) getJWKS(realm string) (*keyfunc.JWKS, error) {
	ja.jwksCacheLock.RLock()
	jwks := ja.jwksCache[realm]
	ja.jwksCacheLock.RUnlock()
	if jwks != nil {
		return jwks, nil
	}

	url := fmt.Sprintf("%s%s/protocol/openid-connect/certs", ja.KeycloakURL, realm)
	log.Printf("[JWTAuth] fetching JWKS from URL=%s", url)
	jwksNew, err := keyfunc.Get(url, ja.jwksOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to load JWKS for realm %q: %w", realm, err)
	}
	ja.jwksCacheLock.Lock()
	ja.jwksCache[realm] = jwksNew
	ja.jwksCacheLock.Unlock()
	return jwksNew, nil
}

func (ja *JWTAuth) verifyAud(claims jwt.MapClaims) error {
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
