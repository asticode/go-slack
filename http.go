package slack

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Send sends an http request with a timeout
var Send = func(req *http.Request, httpClient *http.Client) (*http.Response, error) {
	return httpClient.Do(req)
}

// Send sends a new authorized OHE request
func (o *Slack) Send(hostname string, pattern string, method string, body []byte) (req *http.Request, resp *http.Response, err error) {
	// Log
	url := hostname + pattern
	o.Logger.Debugf("Sending Slack %s request to %s with body %s", method, url, string(body))

	// Create request
	req, err = http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return
	}
	req.Header.Add("Content-type", "application/json")

	// Send request
	resp, err = Send(req, o.HTTPClient)
	return
}

// SendWithMaxRetries sends a new authorized OHE request and retries in case of specific conditions
func (o *Slack) SendWithMaxRetries(hostname string, pattern string, method string, body []byte, retryMax int, retrySleep time.Duration) (req *http.Request, resp *http.Response, err error) {
	// Loop
	for retriesLeft := retryMax; retriesLeft > 0; retriesLeft-- {
		// Send request
		req, resp, err = o.Send(hostname, pattern, method, body)
		if err != nil {
			return
		}

		// Retry if internal server or if too many requests
		if resp.StatusCode >= http.StatusInternalServerError || resp.StatusCode == http.StatusTooManyRequests {
			// Get body
			b, e := ioutil.ReadAll(resp.Body)
			if e != nil {
				err = e
				return
			}

			// Log
			o.Logger.Debugf("Status code %d triggered a retry, sleeping %s and retrying... (%d retries left and body %s)", resp.StatusCode, retrySleep, retriesLeft-1, string(b))

			// Close response body
			resp.Body.Close()

			// Sleep
			time.Sleep(retrySleep)
			continue
		}

		// Return if conditions for retrying were not met
		return
	}

	// Max retries limit reached
	err = fmt.Errorf("Max retries %d reached", retryMax)
	return
}

// ProcessResponse processes an HTTP response
var ProcessResponse = func(req *http.Request, resp *http.Response) error {
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("Invalid status code %v on %v", resp.StatusCode, req.URL.Path)
	}
	return nil
}
