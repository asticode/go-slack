package slack

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Message represents a message
type Message struct {
	Attachments []Attachment `json:"attachments,omitempty"`
	Channel     string       `json:"channel"`
	Markdown    bool         `json:"mrkdwn,omitempty"`
	Text        string       `json:"text,omitempty"`
	Username    string       `json:"username,omitempty"`
}

// Attachment represents an attachments
// https://api.slack.com/docs/message-attachments
type Attachment struct {
	AuthorIcon string  `json:"author_icon,omitempty"`
	AuthorLink string  `json:"author_link,omitempty"`
	AuthorName string  `json:"author_name,omitempty"`
	Color      string  `json:"color,omitempty"`
	Fallback   string  `json:"fallback,omitempty"`
	Fields     []Field `json:"fields,omitempty"`
	Footer     string  `json:"footer,omitempty"`
	FooterIcon string  `json:"footer_icon,omitempty"`
	ImageURL   string  `json:"image_url,omitempty"`
	Pretext    string  `json:"pretext,omitempty"`
	Text       string  `json:"text,omitempty"`
	ThumbURL   string  `json:"thumb_url,omitempty"`
	Title      string  `json:"title,omitempty"`
	TitleLink  string  `json:"title_link,omitempty"`
	Timestamp  int64   `json:"ts,omitempty"`
}

// Field represents an attachment's field
type Field struct {
	Short bool   `json:"short,omitempty"`
	Title string `json:"title,omitempty"`
	Value string `json:"value,omitempty"`
}

// Slack sends a message to Slack
func (s *Slack) Slack(m Message) (err error) {
	// Log
	l := fmt.Sprintf("Slacking message to %s", m.Channel)
	s.Logger.Debugf("[Start] %s", l)
	defer func(now time.Time) {
		s.Logger.Debugf("[End] %s in %s", l, time.Since(now))
	}(time.Now())

	// TODO Make sure texts are HTML encoded

	// Encode message
	var b []byte
	if b, err = json.Marshal(m); err != nil {
		return
	}

	// Send request
	req, resp, err := s.SendWithMaxRetries(s.IncomingWebhookURL, "", http.MethodPost, b)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// Process response
	err = ProcessResponse(req, resp)
	return
}
