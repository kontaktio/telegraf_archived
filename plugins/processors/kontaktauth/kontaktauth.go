package kontaktauth

import (
	_ "embed"
	"encoding/json"
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/processors"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	circuit "github.com/rubyist/circuitbreaker"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type KontaktAuth struct {
	KeycloakURL string `toml:"keycloak_url"`
	ApiAddress  string `toml:"api_address"`
	Audience    string `toml:"audience"`
	CacheSize   int    `toml:"cache_size"`
	ApiCaller   ApiCaller
	JWTAuth     *JWTAuth
}

type apiCompany struct {
	CompanyID string
	Name      string
}

type apiManager struct {
	Company apiCompany
	ID      int64
}

//go:embed sample.conf
var sampleConfig string

var unknownApiKeyDuration = time.Minute * 10

const (
	apiKeyTag    = "Api-Key"
	jwtHeaderTag = "Authorization"
)

var authenticationTime = prometheus.NewHistogram(prometheus.HistogramOpts{
	Name:    "telegraf_authentication_time_seconds",
	Help:    "Time to finish authentication",
	Buckets: prometheus.ExponentialBuckets(0.032, 1.3, 24), // 32ms to ~17s
})

var inFlightRequests = make(map[string]struct{})
var inFlightRequestsLock = sync.Mutex{}

var inFlightRequestsGauge = prometheus.NewGaugeFunc(prometheus.GaugeOpts{
	Name: "telegraf_in_flight_authentication_requests",
	Help: "Number of authentication requests that are currently in flight to Apps-Api",
}, func() float64 {
	return float64(len(inFlightRequests))
})

var cache = make(map[string]apiManager)
var unknownCache = make(map[string]time.Time)

func (ka *KontaktAuth) getManager(apiKey string) (apiManager, error) {
	authenticationStartTime := time.Now()
	defer func() {
		timeDiff := time.Now().Sub(authenticationStartTime).Seconds()
		authenticationTime.Observe(timeDiff)
	}()

	if manager, ok := cache[apiKey]; ok {
		return manager, nil
	}
	if t, ok := unknownCache[apiKey]; ok {
		if t.Add(unknownApiKeyDuration).Before(t) {
			delete(unknownCache, apiKey)
		} else {
			return apiManager{}, errors.New("unauthorized")
		}
	}
	var manager apiManager
	correct, err := ka.get("v2/organization/account/me", apiKey, &manager)
	if err == nil {
		cache[apiKey] = manager
	} else if correct {
		log.Printf("Error %v", err)
		unknownCache[apiKey] = time.Now()
	} else {
		//Don't cache if there wasn't a correct response
		return apiManager{}, err
	}
	return manager, err
}

func (ka *KontaktAuth) get(path, apiKey string, result interface{}) (bool, error) {
	requestId := onRequestStart()
	defer onRequestFinish(requestId)

	response, err := ka.ApiCaller.Call(ka.ApiAddress+path, apiKey)
	if err != nil {
		log.Printf("Error %v", err)
		return false, err
	}
	if response.StatusCode == 401 || response.StatusCode == 403 {
		return true, nil
	}
	if err := json.NewDecoder(response.Body).Decode(result); err != nil {
		return false, err
	}
	return true, nil
}

func onRequestStart() string {
	inFlightRequestsLock.Lock()
	defer inFlightRequestsLock.Unlock()
	requestId := strconv.FormatUint(rand.Uint64(), 16)
	inFlightRequests[requestId] = struct{}{} //Value doesn't matter
	return requestId
}

func onRequestFinish(requestId string) {
	inFlightRequestsLock.Lock()
	defer inFlightRequestsLock.Unlock()
	delete(inFlightRequests, requestId)
}

func (ka *KontaktAuth) SampleConfig() string {
	return sampleConfig
}

func (ka *KontaktAuth) Description() string {
	return "Authenticates telemetry and fills companyId"
}

func (ka *KontaktAuth) Apply(metrics ...telegraf.Metric) []telegraf.Metric {
	result := make([]telegraf.Metric, 0)
	for _, metric := range metrics {
		if metric.HasTag(jwtHeaderTag) {
			tokenStr, _ := metric.GetTag(jwtHeaderTag)
			metric.RemoveTag(jwtHeaderTag)

			claims, err := ka.JWTAuth.VerifyToken(tokenStr)
			if err != nil {
				log.Printf("invalid JWT: %v", err)
				continue
			}

			cid, err := ka.JWTAuth.ExtractCompanyID(claims)
			if err != nil {
				log.Printf("JWT without company-id: %v", err)
				continue
			}

			metric.AddTag("companyId", cid)
			result = append(result, metric)
			continue
		}

		if !metric.HasTag(apiKeyTag) {
			continue
		}
		apiKey, _ := metric.GetTag(apiKeyTag)
		manager, err := ka.getManager(apiKey)
		if err != nil {
			log.Printf("exception while getting manager: %v\n", err)
			continue
		}
		metric.RemoveTag(apiKeyTag)
		metric.AddTag("companyId", manager.Company.CompanyID)
		result = append(result, metric)
	}
	return result
}

func New() *KontaktAuth {
	kontaktAuth := KontaktAuth{
		ApiCaller: &ApiCallerImpl{Client: circuit.NewHTTPClient(5*time.Second, 10, nil)},
	}
	return &kontaktAuth
}

func (ka *KontaktAuth) Init() error {
	KeycloakURL := strings.TrimRight(ka.KeycloakURL, "/") + "/"
	ja := NewJWTAuth(KeycloakURL, ka.Audience, ka.CacheSize)
	ka.JWTAuth = ja
	return nil
}

func init() {
	processors.Add("kontaktauth", func() telegraf.Processor {
		return New()
	})

	prometheus.Register(authenticationTime)
	prometheus.Register(inFlightRequestsGauge)
}

type ApiCaller interface {
	Call(path string, apiKey string) (*http.Response, error)
}

type ApiCallerImpl struct {
	Client *circuit.HTTPClient
}

func (a *ApiCallerImpl) Call(path string, apiKey string) (*http.Response, error) {
	request, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Api-Key", apiKey)
	return a.Client.Do(request)
}
