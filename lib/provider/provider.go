package provider

import (
	"fmt"

	"github.com/mailcow/prometheus-exporter/lib/mailcowApi"
	"github.com/prometheus/client_golang/prometheus"
)

// A Provider is the common abstraction over collection of metrics in this
// exporter. It can provide one or more prometheus collectors (e.g. gauges,
// histograms, ...) that are updated every time the `Update` method is called.
// Be sure to keep a copy of the collectors returned by `GetCollectors`
// in your provider in order to update that same instance.
type Provider interface {
	Name() string
	Provide(mailcowApi.MailcowApiClient) ([]prometheus.Collector, error)
}

func AllProviders() []Provider {
	return []Provider{
		Mailq{},
		Mailbox{},
		Quarantine{},
		Container{},
		Rspamd{},
		Domain{},
		ApiMeta{},
	}
}

func ProviderNames() []string {
	providers := AllProviders()
	names := make([]string, 0, len(providers))
	for _, provider := range providers {
		names = append(names, provider.Name())
	}
	return names
}

func GetProviders(names []string) ([]Provider, error) {
	allProviders := AllProviders()
	providersByName := make(map[string]Provider)
	for _, provider := range allProviders {
		providersByName[provider.Name()] = provider
	}

	selectedProviders := make([]Provider, 0)
	for _, name := range names {
		if provider, exists := providersByName[name]; exists {
			selectedProviders = append(selectedProviders, provider)
		} else {
			return selectedProviders, fmt.Errorf(
				"provider with name '%s' does not exist. Valid provider names are %s",
				name,
				ProviderNames(),
			)
		}
	}

	return selectedProviders, nil
}
