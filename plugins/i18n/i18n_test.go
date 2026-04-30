package i18n

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/Yashwanth-Kumar-26/teleriha/pkg/bot"
)

func TestI18nPlugin(t *testing.T) {
	b := bot.New("test-token")
	
	// Create plugin with some translations
	p := New("en", map[string]map[string]string{
		"en": {
			"welcome": "Welcome %s!",
			"help":    "How can I help you?",
		},
		"es": {
			"welcome": "¡Bienvenido %s!",
			"help":    "¿Cómo puedo ayudarte?",
		},
	})
	
	// Init
	err := p.Init(b)
	assert.NoError(t, err)
	
	// Mock context
	ctx := &bot.Context{
		Sender: &bot.User{ID: 1, LanguageCode: "es"},
	}
	
	// Test translation
	welcome := p.Translate(ctx, "welcome", "Tester")
	assert.Equal(t, "¡Bienvenido Tester!", welcome)
	
	// Test fallback to default language if not specified in user
	ctx2 := &bot.Context{
		Sender: &bot.User{ID: 2}, // No language code
	}
	welcome2 := p.Translate(ctx2, "welcome", "Tester")
	assert.Equal(t, "Welcome Tester!", welcome2)
	
	// Test unknown key
	unknown := p.Translate(ctx, "unknown")
	assert.Equal(t, "unknown", unknown)
}
