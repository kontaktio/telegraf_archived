package kontaktauth

import (
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

func (ja *JWTAuth) VerifyToken(tokenStr string) (string, bool) {
	return ja.Validator.ValidateToken(tokenStr)
}
