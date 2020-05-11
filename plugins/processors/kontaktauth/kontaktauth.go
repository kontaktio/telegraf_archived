package kontaktauth

import (
	"encoding/json"
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/processors"
	"github.com/pkg/errors"
	"github.com/rubyist/circuitbreaker"
	"log"
	"net/http"
	"time"
)

type KontaktAuth struct {
	ApiAddress string `toml:"api_address"`
	ForwardApiKey bool `toml:"forward_api_key"`

	Client *circuit.HTTPClient
}

type apiCompany struct {
	CompanyID string
	Name      string
}

type apiManager struct {
	ID		int64
	EMail   string
	ApiKey  string
	Company apiCompany
}

var SampleConfig = `
api_address="https://testba-api.kontakt.io"
forward_api_key=false
`

var unknownApiKeyDuration = time.Minute * 10

const apiKeyTag = "Api-Key"

var cache = make(map[string]apiManager)
var unknownCache = make(map[string]time.Time)

func (ka *KontaktAuth) getManager(apiKey string) (apiManager, error) {
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
		return apiManager{}, errors.New("error querying manager")
	}
	return manager, err
}

func (ka *KontaktAuth) get(path, apiKey string, result interface{}) (bool, error) {
	request, err := http.NewRequest("GET", ka.ApiAddress+path, nil)
	if err != nil {
		return false, err
	}
	request.Header.Add("Api-Key", apiKey)
	response, err := ka.Client.Do(request)
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

func (p *KontaktAuth) SampleConfig() string {
	return SampleConfig
}

func (p *KontaktAuth) Description() string {
	return "Authenticates telemetry and fills companyId"
}

func (p *KontaktAuth) Apply(metrics ...telegraf.Metric) []telegraf.Metric {
	result := make([]telegraf.Metric, 0)
	for _, metric := range metrics {
		if !metric.HasTag(apiKeyTag) {
			continue
		}
		apiKey, _ := metric.GetTag(apiKeyTag)
		manager, err := p.getManager(apiKey)
		if err != nil {
			continue
		}
		if !p.ForwardApiKey {
			metric.RemoveTag(apiKeyTag)
		}
		metric.AddTag("companyId", manager.Company.CompanyID)
		metric.AddTag("userId", string(manager.ID))
		result = append(result, metric)
	}
	return result
}

func New() *KontaktAuth {
	kontaktAuth := KontaktAuth{
		Client: circuit.NewHTTPClient(time.Second*5, 10, nil),
	}
	return &kontaktAuth
}

func init() {
	processors.Add("kontaktauth", func() telegraf.Processor {
		return New()
	})
}
