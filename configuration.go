package slack

import (
	"flag"
	"time"
)

// Vars
var (
	IncomingWebhookURL = flag.String("slack-incoming-webhook-url", "", "the slack incoming webhook url")
	RequestTimeout     = flag.Duration("slack-request-timeout", 0, "the duration after which a request is considered as having timed out")
)

// Configuration represents the slack configuration
type Configuration struct {
	IncomingWebhookURL string        `toml:"incoming_webhook_url"`
	RequestTimeout     time.Duration `toml:"request_timeout"`
}

// FlagConfig generates a Configuration based on flags
func FlagConfig() Configuration {
	return Configuration{
		IncomingWebhookURL: *IncomingWebhookURL,
		RequestTimeout:     *RequestTimeout,
	}
}
