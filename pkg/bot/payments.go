package bot

// LabeledPrice represents a portion of the price for a product.
type LabeledPrice struct {
	Label  string `json:"label"`
	Amount int    `json:"amount"` // Price in the smallest units of the currency (e.g., cents)
}

// SendInvoiceOption configures a SendInvoice request.
type SendInvoiceOption func(map[string]interface{})

// WithMaxTipAmount sets the maximum accepted amount for tips in the smallest units of the currency.
func WithMaxTipAmount(maxTipAmount int) SendInvoiceOption {
	return func(params map[string]interface{}) {
		params["max_tip_amount"] = maxTipAmount
	}
}

// WithSuggestedTipAmounts sets an array of suggested amounts of tips in the smallest units of the currency.
func WithSuggestedTipAmounts(suggestedTipAmounts []int) SendInvoiceOption {
	return func(params map[string]interface{}) {
		params["suggested_tip_amounts"] = suggestedTipAmounts
	}
}

// WithProviderData sets JSON-encoded data about the invoice.
func WithProviderData(providerData string) SendInvoiceOption {
	return func(params map[string]interface{}) {
		params["provider_data"] = providerData
	}
}

// WithInvoicePhotoURL sets the URL of the product photo for the invoice.
func WithInvoicePhotoURL(photoURL string) SendInvoiceOption {
	return func(params map[string]interface{}) {
		params["photo_url"] = photoURL
	}
}

// SendInvoice sends an invoice.
func (b *Bot) SendInvoice(chatID int64, title, description, payload, providerToken, currency string, prices []LabeledPrice, opts ...SendInvoiceOption) (*Message, error) {
	params := map[string]interface{}{
		"chat_id":     chatID,
		"title":       title,
		"description": description,
		"payload":     payload,
		"currency":    currency,
		"prices":      prices,
	}
	if providerToken != "" {
		params["provider_token"] = providerToken
	}

	for _, opt := range opts {
		opt(params)
	}

	var msg Message
	_, err := b.callMethod("sendInvoice", params, &msg)
	return &msg, err
}

// AnswerPreCheckoutQuery answers a pre-checkout query.
func (b *Bot) AnswerPreCheckoutQuery(preCheckoutQueryID string, ok bool, errorMessage string) error {
	params := map[string]interface{}{
		"pre_checkout_query_id": preCheckoutQueryID,
		"ok":                    ok,
	}
	if !ok {
		params["error_message"] = errorMessage
	}
	_, err := b.callMethod("answerPreCheckoutQuery", params, nil)
	return err
}

// CreateInvoiceLink creates a link for an invoice.
func (b *Bot) CreateInvoiceLink(title, description, payload, providerToken, currency string, prices []LabeledPrice, opts ...SendInvoiceOption) (string, error) {
	params := map[string]interface{}{
		"title":       title,
		"description": description,
		"payload":     payload,
		"currency":    currency,
		"prices":      prices,
	}
	if providerToken != "" {
		params["provider_token"] = providerToken
	}

	for _, opt := range opts {
		opt(params)
	}

	var link string
	_, err := b.callMethod("createInvoiceLink", params, &link)
	return link, err
}

// AnswerShippingQuery answers a shipping query.
func (b *Bot) AnswerShippingQuery(shippingQueryID string, ok bool, shippingOptions []interface{}, errorMessage string) error {
	params := map[string]interface{}{
		"shipping_query_id": shippingQueryID,
		"ok":                ok,
	}

	if ok {
		params["shipping_options"] = shippingOptions
	} else if errorMessage != "" {
		params["error_message"] = errorMessage
	}

	_, err := b.callMethod("answerShippingQuery", params, nil)
	return err
}

// RefundStarPayment refunds a successful payment in Telegram Stars.
func (b *Bot) RefundStarPayment(userID int64, telegramPaymentChargeID string) error {
	params := map[string]interface{}{
		"user_id": userID,
		"telegram_payment_charge_id": telegramPaymentChargeID,
	}
	_, err := b.callMethod("refundStarPayment", params, nil)
	return err
}
