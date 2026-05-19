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

	// Build the Application
	app := bot.NewApplicationBuilder().
		Token(token).
		Opts(bot.WithLogger(log.Logger)).
		Build()

	// Add middleware
	app.Bot.Router.Use(bot.Logger(log.Logger))
	app.Bot.Router.Use(bot.Recover(log.Logger))

	// Register command handlers
	app.AddHandler(bot.CommandHandler("/start", handleStart))
	app.AddHandler(bot.CommandHandler("/help", handleHelp))
	app.AddHandler(bot.CommandHandler("/echo", handleEcho))

	// Register a default handler
	app.Bot.Router.Default(handleDefault)

	// Start the bot in polling mode
	log.Info().Msg("Starting starter bot in polling mode...")
	if err := app.RunPolling(); err != nil {
		log.Fatal().Err(err).Msg("Failed to start bot")
	}

	// Wait for interrupt signal (handled internally by RunPolling)
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

	message := fmt.Sprintf("Hello %s!\n\nWelcome to TeleRiHa!\nUse /help to see what I can do.", name)

	// Send message with a simple inline keyboard
	kb := bot.NewInlineKeyboard(
		bot.NewInlineKeyboardRow(
			bot.InlineButtonURL("Documentation", "https://github.com/Yashwanth-Kumar-26/teleriha"),
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

	_, err := ctx.Reply(fmt.Sprintf("Echo: %s", args))
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

	_, err := ctx.Reply(fmt.Sprintf("%s, %s!", greeting, name))
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
