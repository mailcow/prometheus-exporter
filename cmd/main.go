package main

import (
	_ "embed"
	"fmt"
	"github.com/mailcow/prometheus-exporter/lib/config"
	"github.com/mailcow/prometheus-exporter/lib/mailcowApi"
	"github.com/mailcow/prometheus-exporter/lib/provider"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

func collectMetrics(providers []provider.Provider, conf config.Config) *prometheus.Registry {
	apiClient := mailcowApi.NewMailcowApiClient(
		conf[config.Scheme],
		conf[config.Host],
		conf[config.ApiKey],
	)

	success := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:        "mailcow_exporter_success",
		ConstLabels: map[string]string{"host": conf[config.Host]},
	}, []string{"provider"})

	registry := prometheus.NewRegistry()
	registry.Register(success)

	for _, p := range providers {
		providerSuccess := true
		collectors, err := p.Provide(apiClient)
		if err != nil {
			providerSuccess = false
			log.Printf(
				"Error while updating metrics of %T:\n%s",
				p,
				err.Error(),
			)
		}

		for _, collector := range collectors {
			err = registry.Register(collector)
			if err != nil {
				providerSuccess = false
				log.Printf(
					"Error while updating metrics of %T:\n%s",
					p,
					err.Error(),
				)
			}
		}

		if providerSuccess {
			success.WithLabelValues(fmt.Sprintf("%T", p)).Set(1.0)
		} else {
			success.WithLabelValues(fmt.Sprintf("%T", p)).Set(0.0)
		}
	}

	return registry
}

func main() {
	conf, confSource := config.GetConfig()
	providers := provider.DefaultProviders()

	printConfig(conf, confSource, providers)

	http.HandleFunc("/metrics", func(response http.ResponseWriter, request *http.Request) {
		registry := collectMetrics(providers, conf)

		promhttp.HandlerFor(
			registry,
			promhttp.HandlerOpts{},
		).ServeHTTP(response, request)
	})

	log.Printf("Starting to listen on %s", conf[config.Listen])
	log.Fatal(http.ListenAndServe(conf[config.Listen], nil))
}

func printConfig(config config.Config, confSource config.ConfigSource, providers []provider.Provider) {
	log.Printf("Starting with configuration:")
	for key, value := range config {
		log.Printf("\t%s:\t\"%s\"", key, value)
		log.Printf("\t\t ↳ %s", confSource[key])
	}

	log.Printf("Providers:")
	for _, p := range providers {
		log.Printf("\t%T", p)
	}
}
