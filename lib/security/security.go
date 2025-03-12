package security

import (
	"github.com/mailcow/prometheus-exporter/lib/config"
	"log"
	"net/http"
)

type SecurityProviderCheckResult struct {
	Success         bool
	ExternalMessage string
	InternalMessage string
}

type SecurityProvider interface {
	Check(request http.Request) SecurityProviderCheckResult
	Usage() []string
}

func GetSecurityProvider(conf config.Config) SecurityProvider {
	if conf[config.SecurityInsecureDisableAccessProtection] == "1" {
		log.Printf("[WARNING]")
		log.Printf("[WARNING] Access protection is disabled. This may expose private information to")
		log.Printf("[WARNING] unwanted access. Be sure to add access protection in front of the exporter")
		log.Printf("[WARNING] or only use it internally without access protection.")
		log.Printf("[WARNING]")
		return NoopProvider{}
	}

	token := conf[config.SecurityToken]
	if token == "" {
		log.Printf("[INFO]")
		log.Printf("[INFO] Reusing the Mailcow API key for exporter access. It is recommendde to provide")
		log.Printf("[INFO] a separate key through --security-token / MAILCOW_EXPORTER_SECURITY_TOKEN")
		log.Printf("[INFO]")
		token = conf[config.ApiKey]
	}

	return NewTokenProvider(token)
}
