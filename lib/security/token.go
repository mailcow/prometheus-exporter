package security

import (
	"fmt"
	"net/http"
	"strings"
)

type TokenProvider struct {
	token string
}

func NewTokenProvider(token string) *TokenProvider {
	return &TokenProvider{token: token}
}

func (p TokenProvider) Check(request http.Request) SecurityProviderCheckResult {
	source := "?token="
	token := request.URL.Query().Get("token")

	if token == "" {
		token = extractBearerToken(request)
		source = "Authorization: Bearer "
	}

	if token == "" {
		return SecurityProviderCheckResult{
			Success:         false,
			InternalMessage: "No token provided through query parameter 'token' or Authorization header",
			ExternalMessage: "No token provided",
		}
	}

	if token == p.token {
		return SecurityProviderCheckResult{Success: true}
	}

	return SecurityProviderCheckResult{
		Success:         false,
		InternalMessage: fmt.Sprintf("The token provided through %s%s does not match %s", source, token, p.token),
		ExternalMessage: fmt.Sprintf("The token provided through %s%s is invalid", source, token),
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

func extractBearerToken(request http.Request) string {
	headerValue := request.Header.Get("authorization")
	if headerValue != "" {
		parsed, res := strings.CutPrefix(headerValue, "Bearer ")
		if res {
			return parsed
		}
	}

	return ""
}
