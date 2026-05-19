package gateway

// SendAbilityRequest represents a request to check send ability.
type SendAbilityRequest struct {
	PhoneNumber string `json:"phone_number"`
}

// SendAbilityResponse represents the response from checkSendAbility.
type SendAbilityResponse struct {
	Ok     bool `json:"ok"`
	Result struct {
		RequestID string `json:"request_id"`
	} `json:"result,omitempty"`
	Error string `json:"error,omitempty"`
}

// SendVerificationRequest represents a request to send a verification code.
type SendVerificationRequest struct {
	PhoneNumber   string `json:"phone_number"`
	RequestID     string `json:"request_id,omitempty"`
	SenderUsername string `json:"sender_username,omitempty"`
	CodeLength    int    `json:"code_length,omitempty"`
	Code          string `json:"code,omitempty"`
	TTL           int    `json:"ttl,omitempty"`
	CallbackURL   string `json:"callback_url,omitempty"`
}

// SendVerificationResponse represents the response from sendVerificationMessage.
type SendVerificationResponse struct {
	Ok     bool `json:"ok"`
	Result struct {
		RequestID      string `json:"request_id"`
		DeliveryStatus string `json:"delivery_status"`
	} `json:"result,omitempty"`
	Error string `json:"error,omitempty"`
}

// VerificationStatusRequest represents a request to check the status.
type VerificationStatusRequest struct {
	RequestID string `json:"request_id"`
	Code      string `json:"code,omitempty"`
}

// VerificationStatusResponse represents the response from checkVerificationStatus.
type VerificationStatusResponse struct {
	Ok     bool `json:"ok"`
	Result struct {
		DeliveryStatus string `json:"delivery_status"`
		VerificationStatus string `json:"verification_status,omitempty"`
	} `json:"result,omitempty"`
	Error string `json:"error,omitempty"`
}
