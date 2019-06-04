package slack

import (
	"flag"
	"time"
)

// Vars
var (
	ChannelPrefix      = flag.String("slack-channel-prefix", "", "the slack channel prefix")
	IncomingWebhookURL = flag.String("slack-incoming-webhook-url", "", "the slack incoming webhook url")
	RequestTimeout     = flag.Duration("slack-request-timeout", 0, "the duration after which a request is considered as having timed out")
	RetryMax           = flag.Int("slack-retry-max", 0, "the slack max retry")
	RetrySleep         = flag.Duration("slack-retry-sleep", 0, "the slack max sleep")
)

// Configuration represents the slack configuration
type Configuration struct {
	ChannelPrefix      string        `toml:"channel_prefix"`
	IncomingWebhookURL string        `toml:"incoming_webhook_url"`
	RequestTimeout     time.Duration `toml:"request_timeout"`
	RetryMax           int           `toml:"retry_max"`
	RetrySleep         time.Duration `toml:"retry_sleep"`
}

// FlagConfig generates a Configuration based on flags
func FlagConfig() Configuration {
	return Configuration{
		ChannelPrefix:      *ChannelPrefix,
		IncomingWebhookURL: *IncomingWebhookURL,
		RequestTimeout:     *RequestTimeout,
		RetryMax:           *RetryMax,
		RetrySleep:         *RetrySleep,
	}
}
