package mailcowApi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type FloatStat struct {
	StatusCode int
	Value      float64
}

type BoolStat struct {
	Value bool
}

// Client for mailcow API
type MailcowApiClient struct {
	Scheme        string
	Host          string
	ApiKey        string
	ResponseTimes map[string]FloatStat
	ResponseSizes map[string]FloatStat
	Success       map[string]BoolStat
}

func NewMailcowApiClient(scheme string, host string, apiKey string) MailcowApiClient {
	return MailcowApiClient{
		Scheme:        scheme,
		Host:          host,
		ApiKey:        apiKey,
		ResponseSizes: map[string]FloatStat{},
		ResponseTimes: map[string]FloatStat{},
		Success:       map[string]BoolStat{},
	}
}

// Given an endpoint, this method will do the HTTP request
// with the correct authentication and unserialize the JSON
// response into a given target reference.
func (api MailcowApiClient) Get(endpoint string, target interface{}) error {
	url := fmt.Sprintf("%s://%s/%s", api.Scheme, api.Host, endpoint)
	log.Print(url)

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		api.Success[endpoint] = BoolStat{Value: false}
		return fmt.Errorf(
			"Could not prepare API request to `%s`: %#v",
			endpoint,
			err.Error(),
		)
	}

	request.Header.Add("X-Api-Key", api.ApiKey)
	start := time.Now()

	// API Request
	response, err := (&http.Client{}).Do(request)
	if err != nil {
		api.Success[endpoint] = BoolStat{Value: false}
		return fmt.Errorf(
			"could not execute API request to `%s`: %#v",
			endpoint,
			err.Error(),
		)
	}

	// Metric collection about the API request
	api.ResponseTimes[endpoint] = FloatStat{
		StatusCode: response.StatusCode,
		Value:      float64(time.Since(start).Milliseconds()),
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		api.Success[endpoint] = BoolStat{Value: false}
		return fmt.Errorf(
			"Could not read API response body from endpoint `%s`: \n%s",
			endpoint,
			err.Error(),
		)
	}

	api.ResponseSizes[endpoint] = FloatStat{
		StatusCode: response.StatusCode,
		Value:      float64(len(body)),
	}

	if response.StatusCode != 200 {
		api.Success[endpoint] = BoolStat{Value: false}
		return fmt.Errorf(
			"Received %d response from endpoint `%s`: \n\nResponse body received: \n%s",
			response.StatusCode,
			endpoint,
			body,
		)
	}

	err = json.Unmarshal(body, target)
	if err != nil {
		api.Success[endpoint] = BoolStat{Value: false}
		return fmt.Errorf(
			"Could not parse JSON response from endpoint `%s`: \n%s \n\nResponse body received: \n%s",
			endpoint,
			err.Error(),
			body,
		)
	}

	api.Success[endpoint] = BoolStat{Value: true}
	return nil
}
