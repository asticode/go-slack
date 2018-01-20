package slack

import (
	"net/http"
	"time"
)

// Constants
const (
	ColorDanger  = "danger"
	ColorGood    = "good"
	ColorWarning = "warning"
)

// Slack represents a Slack communicator
type Slack struct {
	ChannelPrefix      string
	HTTPClient         *http.Client
	IncomingWebhookURL string
	RetryMax           int
	RetrySleep         time.Duration
}

// New creates a new Slack communicator
func New(c Configuration) *Slack {
	o := &Slack{
		ChannelPrefix: c.ChannelPrefix,
		HTTPClient: &http.Client{
			Timeout: c.RequestTimeout,
		},
		IncomingWebhookURL: c.IncomingWebhookURL,
		RetryMax:           c.RetryMax,
		RetrySleep:         c.RetrySleep,
	}
	if c.RequestTimeout == 0 {
		o.HTTPClient.Timeout = time.Duration(10) * time.Second
	}
	return o
}
