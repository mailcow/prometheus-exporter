package main

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/mailcow/prometheus-exporter/lib/config"
	"github.com/mailcow/prometheus-exporter/lib/mailcowApi"
	"github.com/mailcow/prometheus-exporter/lib/provider"
	"github.com/mailcow/prometheus-exporter/lib/security"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	conf, confSource, confItems := config.GetConfig()
	securityProvider := security.GetSecurityProvider(conf)
	providers, err := provider.GetProviders(strings.Split(conf[config.Providers], ","))
	if err != nil {
		log.Fatalf("Error while initializing providers: %s", err.Error())
	}

	printConfig(conf, confSource, confItems, providers, securityProvider)

	http.HandleFunc("/metrics", func(response http.ResponseWriter, request *http.Request) {
		checkResult := securityProvider.Check(*request)
		if !checkResult.Success {
			log.Printf("[ERROR]")
			log.Printf("[ERROR] Security provider %T failed request %s %s", securityProvider, request.Method, request.URL)
			log.Printf("[ERROR] External Message: %s", checkResult.ExternalMessage)
			log.Printf("[ERROR] Internal Message: %s", checkResult.InternalMessage)
			log.Printf("[ERROR]")

			response.WriteHeader(http.StatusForbidden)
			response.Write([]byte(checkResult.ExternalMessage))
			return
		}
		registry := collectMetrics(providers, conf)

		promhttp.HandlerFor(
			registry,
			promhttp.HandlerOpts{},
		).ServeHTTP(response, request)
	})

	log.Printf("Starting to listen on %s", conf[config.Listen])
	log.Fatal(http.ListenAndServe(conf[config.Listen], nil))
}

func printConfig(
	config config.Config,
	confSource config.ConfigSource,
	confItems map[config.ConfigKey]config.ConfigItem,
	providers []provider.Provider,
	securityProvider security.SecurityProvider,
) {
	log.Printf("\n")
	log.Printf("Starting with configuration:")
	for key, value := range config {
		log.Printf("\t%s:\t\"%s\"", key, value)
		log.Printf("\t\t ↳ %s", confSource[key])
		log.Printf("\t\t %s", confItems[key].Help)
		log.Printf("\t\t ↳ ENV: %s", confItems[key].EnvVar)
		log.Printf("\t\t ↳ CLI: %s", confItems[key].CliFlag)
		log.Printf("\n")
	}

	log.Printf("\n")
	log.Printf("Security: %T", securityProvider)
	for _, line := range securityProvider.Usage() {
		log.Printf("\t%s", line)
	}
	log.Printf("\n")
	log.Printf("Providers:")
	for _, p := range providers {
		log.Printf("\t%s", p.Name())
	}

	log.Printf("\n")
}
