package slack

import (
	"net/http"
	"time"

	"github.com/rs/xlog"
)

// Constants
const (
	ColorDanger  = "danger"
	ColorGood    = "good"
	ColorWarning = "warning"
)

// Slack represents a Slack communicator
type Slack struct {
	HTTPClient         *http.Client
	IncomingWebhookURL string
	Logger             xlog.Logger
	RetryMax           int
	RetrySleep         time.Duration
}

// New creates a new Slack communicator
func New(c Configuration) *Slack {
	o := &Slack{
		HTTPClient: &http.Client{
			Timeout: c.RequestTimeout,
		},
		IncomingWebhookURL: c.IncomingWebhookURL,
		Logger:             xlog.NopLogger,
		RetryMax:           c.RetryMax,
		RetrySleep:         c.RetrySleep,
	}
	if c.RequestTimeout == 0 {
		o.HTTPClient.Timeout = time.Duration(10) * time.Second
	}
	return o
}
