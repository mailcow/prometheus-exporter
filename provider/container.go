package provider

import (
	"time"

	"github.com/mailcow/prometheus-exporter/mailcowApi"
	"github.com/prometheus/client_golang/prometheus"
)

type Container struct{}

type containerItem struct {
	Container string `json:"container"`
	State     string `json:"state"`
	StartedAt string `json:"started_at"`
	Image     string `json:"image"`
}

func (container Container) Provide(api mailcowApi.MailcowApiClient) ([]prometheus.Collector, error) {
	startTime := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:        "mailcow_container_start",
		Help:        "Unix timestamp of the container start",
		ConstLabels: map[string]string{"host": api.Host},
	}, []string{"container", "image"})
	running := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:        "mailcow_container_running",
		Help:        "1 if the container is running, 0 if not",
		ConstLabels: map[string]string{"host": api.Host},
	}, []string{"container", "image"})
	collectors := []prometheus.Collector{running, startTime}

	body := make(map[string]containerItem)
	err := api.Get("api/v1/get/status/containers", &body)
	if err != nil {
		return collectors, err
	}

	for _, item := range body {
		isRunning := 0.0
		if item.State == "running" {
			isRunning = 1.0
		}

		t, err := time.Parse(time.RFC3339Nano, item.StartedAt)
		if err != nil {
			return []prometheus.Collector{}, err
		}

		running.WithLabelValues(item.Container, item.Image).Set(isRunning)
		startTime.WithLabelValues(item.Container, item.Image).Set(float64(t.Unix()))
	}

	return collectors, nil
}
