package slack_test

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/asticode/go-slack"
	"github.com/rs/xlog"
	"github.com/stretchr/testify/assert"
)

func TestSendWithMaxRetries(t *testing.T) {
	var count int
	o := slack.Slack{
		Logger: xlog.NopLogger,
	}
	slack.Send = func(req *http.Request, httpClient *http.Client) (*http.Response, error) {
		count++
		if count == 1 {
			return &http.Response{StatusCode: http.StatusInternalServerError, ProtoMinor: 1, Body: ioutil.NopCloser(strings.NewReader(""))}, nil
		} else if count == 2 {
			return &http.Response{StatusCode: http.StatusGatewayTimeout, ProtoMinor: 2, Body: ioutil.NopCloser(strings.NewReader(""))}, nil
		} else if count == 3 {
			return &http.Response{StatusCode: http.StatusTooManyRequests, ProtoMinor: 3, Body: ioutil.NopCloser(strings.NewReader(""))}, nil
		}
		return &http.Response{StatusCode: http.StatusBadRequest, ProtoMinor: 4, Body: ioutil.NopCloser(strings.NewReader(""))}, nil
	}
	_, resp, err := o.SendWithMaxRetries("", "", "", nil, 1, time.Nanosecond)
	assert.Error(t, err)
	assert.Equal(t, 1, resp.ProtoMinor)
	count = 0
	_, resp, err = o.SendWithMaxRetries("", "", "", nil, 2, time.Nanosecond)
	assert.Error(t, err)
	assert.Equal(t, 2, resp.ProtoMinor)
	count = 0
	_, resp, err = o.SendWithMaxRetries("", "", "", nil, 3, time.Nanosecond)
	assert.Error(t, err)
	assert.Equal(t, 3, resp.ProtoMinor)
	count = 0
	_, resp, err = o.SendWithMaxRetries("", "", "", nil, 4, time.Nanosecond)
	assert.NoError(t, err)
	assert.Equal(t, 4, resp.ProtoMinor)
}

func TestProcessResponse(t *testing.T) {
	assert.Error(t, slack.ProcessResponse(&http.Request{URL: &url.URL{RawQuery: ""}}, &http.Response{StatusCode: http.StatusBadRequest, Body: ioutil.NopCloser(strings.NewReader(""))}))
	assert.Error(t, slack.ProcessResponse(&http.Request{URL: &url.URL{RawQuery: ""}}, &http.Response{StatusCode: http.StatusInternalServerError, Body: ioutil.NopCloser(strings.NewReader(""))}))
	assert.NoError(t, slack.ProcessResponse(&http.Request{URL: &url.URL{RawQuery: ""}}, &http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(strings.NewReader(""))}))
	assert.NoError(t, slack.ProcessResponse(&http.Request{URL: &url.URL{RawQuery: ""}}, &http.Response{StatusCode: http.StatusCreated, Body: ioutil.NopCloser(strings.NewReader(""))}))
	assert.NoError(t, slack.ProcessResponse(&http.Request{URL: &url.URL{RawQuery: ""}}, &http.Response{StatusCode: http.StatusNoContent, Body: ioutil.NopCloser(strings.NewReader(""))}))
}
