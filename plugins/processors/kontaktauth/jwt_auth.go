package kontaktauth

import (
	"errors"
	"log"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/coocood/freecache"
	"github.com/golang-jwt/jwt/v4"
)

const (
	defaultAudience  = "compute-api"
	defaultCacheSize = 1 << 24
)

type JWTAuth struct {
	Validator TokenValidator
}

func NewJWTAuth(KeycloakURL string, Audience string, CacheSize int) *JWTAuth {
	if Audience == "" {
		Audience = defaultAudience
	}
	if CacheSize <= 0 {
		CacheSize = defaultCacheSize
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
	cache := freecache.NewCache(CacheSize)
	caching := &CachingValidator{
		base:       base,
		tokenCache: cache,
	}
	return &JWTAuth{Validator: caching}
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
	return ja.Validator.ValidateToken(tokenStr)
}
