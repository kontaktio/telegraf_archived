package kontaktauth

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
)

const (
	companyIdClaim  = "company-id"
	defaultAudience = "compute-api"
)

type JWTAuth struct {
	Validator TokenValidator
}

func NewJWTAuth(KeycloakURL string, Audience string) *JWTAuth {
	if Audience == "" {
		Audience = defaultAudience
	}

	base := &JwksValidator{
		KeycloakURL: KeycloakURL,
		Audience:    Audience,
		jwksCache:   make(map[string]*keyfunc.JWKS),
		jwksOpts: keyfunc.Options{
			RefreshErrorHandler: func(err error) {
				log.Printf("JWKS refresh error: %v", err)
			},
			RefreshInterval:   time.Hour,
			RefreshUnknownKID: true,
		},
	}
	caching := &CachingValidator{
		base:      base,
		cache:     make(map[string][]*cacheEntry),
		jwtParser: new(jwt.Parser),
	}
	return &JWTAuth{Validator: caching}
}

func ExtractCompanyID(tokenStr string) (string, error) {
	parser := new(jwt.Parser)
	token, _, err := parser.ParseUnverified(tokenStr, jwt.MapClaims{})
	if err != nil {
		return "", err
	}

	claims := token.Claims.(jwt.MapClaims)
	val, ok := claims[companyIdClaim]
	if !ok {
		return "", fmt.Errorf("claim %q not found", companyIdClaim)
	}
	companyId, ok := val.(string)
	if !ok {
		return "", errors.New("company-id claim is not a string")
	}
	return companyId, nil
}

func (ja *JWTAuth) VerifyToken(tokenStr string) bool {
	return ja.Validator.ValidateToken(tokenStr)
}
