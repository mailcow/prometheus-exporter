package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type ConfigKey string

const (
	Scheme                                  ConfigKey = "scheme"
	Host                                    ConfigKey = "host"
	ApiKey                                  ConfigKey = "apiKey"
	Listen                                  ConfigKey = "listen"
	SecurityToken                           ConfigKey = "securityToken"
	SecurityInsecureDisableAccessProtection ConfigKey = "securityInsecureDisableAccessProtection"
)

type Config map[ConfigKey]string
type ConfigSource map[ConfigKey]string
type ConfigItem struct {
	EnvVar   string
	Help     string
	CliFlag  string
	Default  string
	Required bool
}

var (
	configItems = map[ConfigKey]ConfigItem{
		Scheme: ConfigItem{
			EnvVar:  "MAILCOW_EXPORTER_SCHEME",
			CliFlag: "scheme",
			Default: "https",
			Help:    "The Scheme used to access the mailcow API (http/https)",
		},
		Host: ConfigItem{
			EnvVar:   "MAILCOW_EXPORTER_HOST",
			CliFlag:  "host",
			Required: true,
			Help:     "The hostname of the mailcow instance",
		},
		ApiKey: ConfigItem{
			EnvVar:   "MAILCOW_EXPORTER_API_KEY",
			CliFlag:  "api-key",
			Required: true,
			Help:     "The API Key used to access mailcow",
		},
		Listen: ConfigItem{
			EnvVar:  "MAILCOW_EXPORTER_LISTEN",
			CliFlag: "listen",
			Default: ":9099",
			Help:    "Hostname and port the exporter listens on",
		},
		SecurityToken: ConfigItem{
			EnvVar:  "MAILCOW_EXPORTER_SECURITY_TOKEN",
			CliFlag: "security-token",
			Help:    "The token that must be provided through the ?token=... URL parameter when accessing the exporter.\nThis defaults to the mailcow API Token if nothing is provided.",
		},
		SecurityInsecureDisableAccessProtection: ConfigItem{
			EnvVar:  "MAILCOW_EXPORTER_SECURITY_DISABLE_ACCESS_PROTECTION",
			CliFlag: "security-insecure-disable-access-protection",
			Default: "0",
			Help:    "Disables access protection for the exporter if set to 1.\nThis may expose sensitive information if no other security precautions are used",
		},
	}
)

func GetConfig() (Config, ConfigSource) {
	// Gather flags
	flagValues := map[ConfigKey]*string{}
	for key, configItem := range configItems {
		flagValues[key] = flag.String(
			configItem.CliFlag,
			"",
			buildHelpString(configItem),
		)
	}

	// Parse Flags
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Precedence of config vars: CLI flag overrides environment variable overrides default value")
		fmt.Fprintf(os.Stderr, "\n\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	config := Config{}
	configSource := ConfigSource{}
	// Apply default values to config
	for key, configItem := range configItems {
		config[key] = configItem.Default
		configSource[key] = "default value"
	}

	// Apply environment variables to config
	for key, configItem := range configItems {
		env := os.Getenv(configItem.EnvVar)
		if env != "" {
			config[key] = env
			configSource[key] = fmt.Sprintf("env %s", configItem.EnvVar)
		}
	}

	// Apply flags to config
	for key, _ := range configItems {
		if *flagValues[key] != "" {
			config[key] = *flagValues[key]
			configSource[key] = fmt.Sprintf("CLI: --%s=%s", key, *flagValues[key])
		}
	}

	validationErrors := []string{}
	for key, configItem := range configItems {
		if configItem.Required && config[key] == "" {
			validationErrors = append(validationErrors, buildValidationError(key, configItem))
		}
	}
	if len(validationErrors) > 0 {
		log.Fatalf("Configuration validation errors:\n\t%s", strings.Join(validationErrors, "\n\t"))
	}

	return config, configSource
}

func buildHelpString(item ConfigItem) string {
	additions := []string{}

	if item.Default != "" {
		additions = append(additions, fmt.Sprintf("default:  \"%s\"", item.Default))
	}
	if item.EnvVar != "" {
		env := os.Getenv(item.EnvVar)
		additions = append(additions, fmt.Sprintf("env:      %s=\"%s\"", item.EnvVar, env))
	}
	additions = append(additions, fmt.Sprintf("CLI:      --%s", item.CliFlag))
	if item.Required {
		additions = append(additions, "[required]")
	}

	help := item.Help
	for _, addition := range additions {
		help += fmt.Sprintf("\n\t%s", addition)
	}

	return help
}

func buildValidationError(key ConfigKey, item ConfigItem) string {
	methods := []string{}
	if item.EnvVar != "" {
		methods = append(methods, fmt.Sprintf("env: %s", item.EnvVar))
	}
	if item.CliFlag != "" {
		methods = append(methods, fmt.Sprintf("cli: --%s", item.CliFlag))
	}

	return fmt.Sprintf("%s: required. Provide it with one of the following methods: %s", key, strings.Join(methods, ", "))
}
