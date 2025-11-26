package provider

import (
	"strconv"

	"github.com/mailcow/prometheus-exporter/lib/mailcowApi"
	"github.com/prometheus/client_golang/prometheus"
)

type ApiMeta struct{}

func (self ApiMeta) Name() string {
	return "ApiMeta"
}

func (self ApiMeta) Provide(api mailcowApi.MailcowApiClient) ([]prometheus.Collector, error) {
	responseTime := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:        "mailcow_api_response_time",
		Help:        "Response time of the API in milliseconds (1/1000s of a second)",
		ConstLabels: map[string]string{"host": api.Host},
	}, []string{"endpoint", "statusCode"})
	responseSize := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:        "mailcow_api_response_size",
		Help:        "Size of API response in bytes",
		ConstLabels: map[string]string{"host": api.Host},
	}, []string{"endpoint", "statusCode"})
	success := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:        "mailcow_api_success",
		Help:        "1, if request was sucessful, 0 if not",
		ConstLabels: map[string]string{"host": api.Host},
	}, []string{"endpoint"})

	collectors := []prometheus.Collector{responseTime, responseSize, success}

	for endpoint, item := range api.ResponseTimes {
		statusCodeString := strconv.FormatInt(int64(item.StatusCode), 10)
		responseTime.WithLabelValues(endpoint, statusCodeString).Set(item.Value)
	}
	for endpoint, item := range api.ResponseSizes {
		statusCodeString := strconv.FormatInt(int64(item.StatusCode), 10)
		responseSize.WithLabelValues(endpoint, statusCodeString).Set(item.Value)
	}
	for endpoint, item := range api.Success {
		if item.Value {
			success.WithLabelValues(endpoint).Set(1.0)
		} else {
			success.WithLabelValues(endpoint).Set(0.0)
		}
	}

	return collectors, nil
}
