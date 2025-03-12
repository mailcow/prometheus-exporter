package provider

import (
	"time"

	"github.com/mailcow/prometheus-exporter/mailcowApi"
	"github.com/prometheus/client_golang/prometheus"
)

// Quarantine Provider. Use `NewQuarantine` to initialize this struct.
// This provider uses the `/api/v1/get/quarantine/all` endpoint
// in order to gather metrics about quarantined mails.
type Quarantine struct{}

type quarantineItem struct {
	VirusFlag int     `json:"virus_flag"`
	Score     float64 `json:"score"`
	Recipient string  `json:"rcpt"`
	Created   int64   `json:"created"`
}

func (quarantine Quarantine) Provide(api mailcowApi.MailcowApiClient) ([]prometheus.Collector, error) {
	countGauge := *prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:        "mailcow_quarantine_count",
		Help:        "Number of mails currently in quarantine",
		ConstLabels: map[string]string{"host": api.Host},
	}, []string{"recipient", "is_virus"})
	scoreHist := *prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:        "mailcow_quarantine_score",
		Help:        "Score of quarantined mails according to rspamd",
		Buckets:     []float64{0, 10, 20, 40, 60, 80, 100},
		ConstLabels: map[string]string{"host": api.Host},
	}, []string{"recipient"})
	ageHist := *prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "mailcow_quarantine_age",
		Help: "Age of quarantined items in seconds",
		Buckets: []float64{
			(60 * 60 * 3),       // 3 hours
			(60 * 60 * 12),      // 12 hours
			(60 * 60 * 24),      // 1 day
			(3 * 60 * 60 * 24),  // 3 days
			(7 * 60 * 60 * 24),  // 7 days
			(14 * 60 * 60 * 24), // 14 days
			(30 * 60 * 60 * 24), // 30 days
		},
		ConstLabels: map[string]string{"host": api.Host},
	}, []string{"recipient"})

	collectors := []prometheus.Collector{
		countGauge,
		scoreHist,
		ageHist,
	}

	body := make([]quarantineItem, 0)
	err := api.Get("api/v1/get/quarantine/all", &body)
	if err != nil {
		return collectors, err
	}

	virus := make(map[string]int)
	notVirus := make(map[string]int)
	for _, q := range body {
		if _, ok := virus[q.Recipient]; !ok {
			virus[q.Recipient] = 0
		}
		if _, ok := notVirus[q.Recipient]; !ok {
			notVirus[q.Recipient] = 0
		}

		if q.VirusFlag == 1 {
			virus[q.Recipient]++
		} else {
			notVirus[q.Recipient]++
		}

		age := time.Now().Unix() - q.Created
		ageHist.WithLabelValues(q.Recipient).Observe(float64(age))
		scoreHist.WithLabelValues(q.Recipient).Observe(float64(q.Score))
	}

	for recipient, count := range virus {
		countGauge.WithLabelValues(recipient, "1").Set(float64(count))
	}
	for recipient, count := range notVirus {
		countGauge.WithLabelValues(recipient, "0").Set(float64(count))
	}

	return collectors, nil
}
