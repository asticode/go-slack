package slack

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/asticode/go-astilog"
)

// Send sends an http request with a timeout
var Send = func(req *http.Request, httpClient *http.Client) (*http.Response, error) {
	return httpClient.Do(req)
}

// Send sends a new authorized OHE request
func (s *Slack) Send(hostname string, pattern string, method string, body []byte) (req *http.Request, resp *http.Response, err error) {
	// Log
	url := hostname + pattern
	astilog.Debugf("Sending Slack %s request to %s with body %s", method, url, string(body))

	// Create request
	req, err = http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return
	}
	req.Header.Add("Content-type", "application/json")
	defer req.Body.Close()

	// Send request
	if resp, err = Send(req, s.HTTPClient); err != nil {
		astilog.Error(fmt.Sprintf("%s for request to %s", err, req.URL))
	}
	return
}

// SendWithMaxRetries sends a new authorized OHE request and retries in case of specific conditions
func (s *Slack) SendWithMaxRetries(hostname string, pattern string, method string, body []byte) (req *http.Request, resp *http.Response, err error) {
	// Loop
	// We start at s.RetryMax + 1 so that it runs at least once even if RetryMax == 0
	for retriesLeft := s.RetryMax + 1; retriesLeft > 0; retriesLeft-- {
		// Send request
		var retry bool
		if req, resp, err = s.Send(hostname, pattern, method, body); err != nil {
			// If error is temporary, retry
			if netError, ok := err.(net.Error); ok && netError.Temporary() {
				retry = true
			} else {
				return
			}
		}

		// Retry if internal server or if too many requests
		if retry || resp.StatusCode >= http.StatusInternalServerError || resp.StatusCode == http.StatusTooManyRequests {
			// Get body
			var b []byte
			if resp != nil {
				defer resp.Body.Close()
				if b, err = ioutil.ReadAll(resp.Body); err != nil {
					astilog.Error(err)
					return
				}
			}

			// Log
			if retriesLeft > 1 {
				astilog.Debugf("Sleeping %s and retrying... (%d retries left and body %s)", s.RetrySleep, retriesLeft-1, string(b))
				time.Sleep(s.RetrySleep)
			}
			continue
		}

		// Return if conditions for retrying were not met
		return
	}

	// Max retries limit reached
	err = fmt.Errorf("Max retries %d reached for request to %s", s.RetryMax, req.URL)
	astilog.Error(err)
	return
}

// ProcessResponse processes an HTTP response
var ProcessResponse = func(req *http.Request, resp *http.Response) error {
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("Invalid status code %v on %v", resp.StatusCode, req.URL)
	}
	return nil
}
