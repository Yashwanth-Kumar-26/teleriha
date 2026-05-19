package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const defaultBaseURL = "https://gatewayapi.telegram.org/"

// Client interacts with the Telegram Gateway API.
type Client struct {
	Token      string
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient creates a new Gateway API client.
func NewClient(token string) *Client {
	return &Client{
		Token:   token,
		BaseURL: defaultBaseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// call makes a POST request to the Gateway API.
func (c *Client) call(ctx context.Context, method string, payload interface{}, out interface{}) error {
	url := c.BaseURL + method

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.Token)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("http request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	if err := json.Unmarshal(respBody, out); err != nil {
		return fmt.Errorf("failed to decode response: %w\n%s", err, string(respBody))
	}

	return nil
}

// CheckSendAbility checks if a message can be sent to a specific phone number.
func (c *Client) CheckSendAbility(ctx context.Context, phoneNumber string) (*SendAbilityResponse, error) {
	req := SendAbilityRequest{PhoneNumber: phoneNumber}
	var res SendAbilityResponse
	err := c.call(ctx, "checkSendAbility", req, &res)
	if err != nil {
		return nil, err
	}
	if !res.Ok {
		return &res, fmt.Errorf("gateway API error: %s", res.Error)
	}
	return &res, nil
}

// SendVerificationMessage sends a verification code to a phone number.
func (c *Client) SendVerificationMessage(ctx context.Context, req SendVerificationRequest) (*SendVerificationResponse, error) {
	var res SendVerificationResponse
	err := c.call(ctx, "sendVerificationMessage", req, &res)
	if err != nil {
		return nil, err
	}
	if !res.Ok {
		return &res, fmt.Errorf("gateway API error: %s", res.Error)
	}
	return &res, nil
}

// CheckVerificationStatus checks the delivery or verification status of a request.
func (c *Client) CheckVerificationStatus(ctx context.Context, requestID, code string) (*VerificationStatusResponse, error) {
	req := VerificationStatusRequest{
		RequestID: requestID,
		Code:      code,
	}
	var res VerificationStatusResponse
	err := c.call(ctx, "checkVerificationStatus", req, &res)
	if err != nil {
		return nil, err
	}
	if !res.Ok {
		return &res, fmt.Errorf("gateway API error: %s", res.Error)
	}
	return &res, nil
}
