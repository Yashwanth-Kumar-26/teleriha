// Package main provides an AI-powered bot example using TeleRiHa.
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/Yashwanth-Kumar-26/teleriha/pkg/bot"
)

// AIProvider is a simple interface for AI responses.
type AIProvider interface {
	GenerateResponse(prompt string) (string, error)
}

// MockAI is a simple mock AI provider.
type MockAI struct{}

func (m *MockAI) GenerateResponse(prompt string) (string, error) {
	// Simple mock logic: echo back with some "AI" flavor
	return fmt.Sprintf("🤖 AI Response to '%s':\n\nI have analyzed your input and concluded that you are asking something interesting! (Mock Response)", prompt), nil
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Debug().Err(err).Msg("No .env file found")
	}

	// Get bot token
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		// Use a dummy token for the example if not provided
		token = "123456789:ABCDEF123456789"
		log.Warn().Msg("BOT_TOKEN not found, using dummy token for example")
	}

	// Configure logging
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Create a new bot
	b := bot.New(token, bot.WithLogger(log.Logger))

	// Add middleware
	b.Router.Use(bot.Logger(log.Logger))
	b.Router.Use(bot.Recover(log.Logger))

	ai := &MockAI{}

	// Register command handlers
	b.Router.On("/start", func(ctx *bot.Context) error {
		_, err := ctx.Reply("Hello! I am an AI-powered bot. Send me any message and I will try to respond!")
		return err
	})

	// Default handler for all other messages (pass to AI)
	b.Router.Default(func(ctx *bot.Context) error {
		text := ctx.Text()
		if text == "" || strings.HasPrefix(text, "/") {
			return nil
		}

		response, err := ai.GenerateResponse(text)
		if err != nil {
			_, replyErr := ctx.Reply("Sorry, I encountered an error while processing your request.")
			return replyErr
		}

		_, replyErr := ctx.Reply(response)
		return replyErr
	})

	// Start the bot (in this example we don't actually block so it can be built)
	log.Info().Msg("AI bot example initialized")
	
	// If running as a real bot, we would use:
	// b.StartPolling(1*time.Second, 30*time.Second)
	
	fmt.Println("This is an example bot. To run it, provide a BOT_TOKEN and uncomment StartPolling.")
}
