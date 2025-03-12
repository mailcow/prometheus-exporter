package provider

import (
	"github.com/mailcow/prometheus-exporter/lib/mailcowApi"
	"github.com/prometheus/client_golang/prometheus"
)

// A Provider is the common abstraction over collection of metrics in this
// exporter. It can provide one or more prometheus collectors (e.g. gauges,
// histograms, ...) that are updated every time the `Update` method is called.
// Be sure to keep a copy of the collectors returned by `GetCollectors`
// in your provider in order to update that same instance.
type Provider interface {
	Provide(mailcowApi.MailcowApiClient) ([]prometheus.Collector, error)
}

func DefaultProviders() []Provider {
	return []Provider{
		ApiMeta{},
		Mailq{},
		Mailbox{},
		Quarantine{},
		Container{},
		Rspamd{},
		Domain{},
	}
}
