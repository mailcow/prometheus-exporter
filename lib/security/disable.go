package security

import "net/http"

type NoopProvider struct{}

func (p NoopProvider) Check(request http.Request) SecurityProviderCheckResult {
	return SecurityProviderCheckResult{Success: true}
}

func (p NoopProvider) Usage() []string {
	return []string{
		"scrape_configs:",
		"  - job_name: 'mailcow'",
		"    static_configs:",
		"      - targets: ['mailcow-exporter-hostname:9099']",
	}
}
