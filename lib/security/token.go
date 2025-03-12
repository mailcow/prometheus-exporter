package security

import (
	"fmt"
	"net/http"
)

type TokenProvider struct {
	token string
}

func NewTokenProvider(token string) *TokenProvider {
	return &TokenProvider{token: token}
}

func (p TokenProvider) Check(request http.Request) SecurityProviderCheckResult {
	token := request.URL.Query().Get("token")
	if token == p.token {
		return SecurityProviderCheckResult{Success: true}
	}

	return SecurityProviderCheckResult{
		Success:         false,
		InternalMessage: fmt.Sprintf("The token provided through ?token=%s does not match %s", token, p.token),
		ExternalMessage: fmt.Sprintf("The token provided through ?token=%s is invalid", token),
	}
}

func (p TokenProvider) Usage() []string {
	return []string{
		"scrape_configs:",
		"  - job_name: 'mailcow'",
		"    static_configs:",
		"      - targets: ['mailcow-exporter-hostname:9099']",
		"    params:",
		fmt.Sprintf("      token: ['%s']", p.token),
	}
}
