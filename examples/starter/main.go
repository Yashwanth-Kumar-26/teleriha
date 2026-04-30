// Package main provides a simple starter bot example using TeleRiHa.
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/Yashwanth-Kumar-26/teleriha/pkg/bot"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Debug().Err(err).Msg("No .env file found")
	}

	// Get bot token
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatal().Msg("BOT_TOKEN environment variable is required")
	}

	// Configure logging
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Create a new bot
	b := bot.New(token, bot.WithLogger(log.Logger))

	// Add middleware
	b.Router.Use(bot.Logger(log.Logger))
	b.Router.Use(bot.Recover(log.Logger))

	// Register command handlers
	b.Router.On("/start", handleStart)
	b.Router.On("/help", handleHelp)
	b.Router.On("/echo", handleEcho)

	// Register a regex handler
	b.Router.OnRegex(`^(?i)hello|hi|hey$`, handleGreeting)

	// Register a default handler
	b.Router.Default(handleDefault)

	// Start the bot in polling mode
	log.Info().Msg("Starting starter bot in polling mode...")
	if err := b.StartPolling(1*time.Second, 30*time.Second); err != nil {
		log.Fatal().Err(err).Msg("Failed to start bot")
	}

	// Wait for interrupt signal (handled internally by StartPolling)
	select {}
}

// handleStart handles the /start command.
func handleStart(ctx *bot.Context) error {
	sender := ctx.Sender
	if sender == nil {
		return nil
	}

	name := sender.FirstName
	if sender.LastName != "" {
		name += " " + sender.LastName
	}

	message := fmt.Sprintf("Hello %s! 👋\n\nWelcome to TeleRiHa!\nUse /help to see what I can do.", name)

	// Send message with a simple inline keyboard
	kb := bot.NewInlineKeyboard(
		bot.NewInlineKeyboardRow(
			bot.InlineButtonURL("📖 Documentation", "https://github.com/Yashwanth-Kumar-26/teleriha"),
		),
	)

	_, err := ctx.Reply(message, bot.WithReplyMarkup(kb))
	return err
}

// handleHelp handles the /help command.
func handleHelp(ctx *bot.Context) error {
	helpText := `Available commands:

/start - Start the bot
/help - Show this help message
/echo <text> - Echo back the provided text

You can also say hello, hi, or hey!

This bot is built with TeleRiHa, a production-grade Telegram bot framework for Go.`

	_, err := ctx.Reply(helpText)
	return err
}

// handleEcho handles the /echo command.
func handleEcho(ctx *bot.Context) error {
	args := ctx.CommandArgs()
	if args == "" {
		_, err := ctx.Reply("Please provide some text to echo. Usage: /echo <text>")
		return err
	}

	_, err := ctx.Reply(fmt.Sprintf("🔁 You said: %s", args))
	return err
}

// handleGreeting handles greetings like hello, hi, hey.
func handleGreeting(ctx *bot.Context) error {
	sender := ctx.Sender
	if sender == nil {
		return nil
	}

	name := sender.FirstName
	if sender.LastName != "" {
		name += " " + sender.LastName
	}

	greetings := []string{"Hello", "Hi", "Hey", "Hi there", "Hello there"}
	greeting := greetings[time.Now().Unix()%int64(len(greetings))]

	_, err := ctx.Reply(fmt.Sprintf("%s, %s! 😊", greeting, name))
	return err
}

// handleDefault handles messages that don't match any command.
func handleDefault(ctx *bot.Context) error {
	text := ctx.Text()
	if text == "" {
		// Not a text message
		return nil
	}

	_, err := ctx.Reply(fmt.Sprintf("I don't understand '%s'. Try /help to see what I can do.", text))
	return err
}
