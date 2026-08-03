// Package notification is an example outbound API client. Every external
// service gets its own package like this one: an interface the service
// layer depends on, plus request/response models that describe *that
// service's* wire format — never reused from dto/, and never leaking
// into it. If the notification provider changes its field names
// tomorrow, only this package changes.
package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Notifier interface {
	SendWelcomeEmail(ctx context.Context, req SendWelcomeEmailRequest) (*SendWelcomeEmailResponse, error)
}

type httpNotifier struct {
	baseURL string
	apiKey  string
	http    *http.Client
}

func NewHTTPNotifier(baseURL, apiKey string) Notifier {
	return &httpNotifier{
		baseURL: baseURL,
		apiKey:  apiKey,
		http:    &http.Client{Timeout: 5 * time.Second},
	}
}

func (c *httpNotifier) SendWelcomeEmail(ctx context.Context, req SendWelcomeEmailRequest) (*SendWelcomeEmailResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v1/emails/welcome", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("notification: unexpected status %d", resp.StatusCode)
	}

	var out SendWelcomeEmailResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}
