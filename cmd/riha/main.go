package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/Yashwanth-Kumar-26/teleriha/pkg/bot"
)

var (
	version = "0.1.0"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Debug().Err(err).Msg("No .env file found")
	}

	// Create root command
	rootCmd := &cobra.Command{
		Use:           "riha",
		Short:         "TeleRiHa - Telegram Bot Framework",
		Long:          "Build Telegram bots at the speed of thought.\n\nTeleRiHa is a production-grade Telegram bot framework for Go.",
		Version:       fmt.Sprintf("%s (commit: %s, date: %s)", version, commit, date),
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Configure logging
	configureLogging(rootCmd)

	// Add commands
	rootCmd.AddCommand(
		newRunCommand(),
		newDevCommand(),
		newVersionCommand(),
		newWebhookCommand(),
	)

	// Execute
	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Err(err).Msg("Failed to execute command")
	}
}

func configureLogging(cmd *cobra.Command) {
	// Set up zerolog
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Check for debug flag
	debug, _ := cmd.Flags().GetBool("debug")
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// Use console output with colors
	log.Logger = zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
	}).With().Timestamp().Logger()
}

func newRunCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run the Telegram bot",
		Long:  "Run the Telegram bot in polling mode.",
		RunE:  runPolling,
	}

	cmd.Flags().String("token", "", "Bot token from @BotFather")
	cmd.Flags().String("env", "", "Path to .env file")
	cmd.Flags().Int("poll-interval", 1, "Polling interval in seconds")
	cmd.Flags().Int("poll-timeout", 30, "Polling timeout in seconds")
	cmd.Flags().Bool("debug", false, "Enable debug logging")

	return cmd
}

func newDevCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dev",
		Short: "Run the Telegram bot in development mode",
		Long:  "Run the Telegram bot in development mode with hot reload.",
		RunE:  runDev,
	}

	cmd.Flags().String("token", "", "Bot token from @BotFather")
	cmd.Flags().String("env", "", "Path to .env file")
	cmd.Flags().Int("port", 8080, "Port for the webhook server")
	cmd.Flags().String("webhook-path", "/webhook", "Path for webhook endpoint")
	cmd.Flags().Bool("debug", false, "Enable debug logging")

	return cmd
}

func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Long:  "Print the version, commit hash, and build date of TeleRiHa.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("TeleRiHa version %s\n", version)
			fmt.Printf("Commit: %s\n", commit)
			fmt.Printf("Built: %s\n", date)
		},
	}
}

func newWebhookCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "webhook",
		Short: "Manage webhook settings",
		Long:  "Manage the bot's webhook settings.",
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:   "set",
			Short: "Set webhook URL",
			RunE:  setWebhook,
		},
		&cobra.Command{
			Use:   "delete",
			Short: "Delete webhook",
			RunE:  deleteWebhook,
		},
		&cobra.Command{
			Use:   "info",
			Short: "Get webhook info",
			RunE:  getWebhookInfo,
		},
	)

	cmd.PersistentFlags().String("token", "", "Bot token from @BotFather")
	cmd.PersistentFlags().String("url", "", "Webhook URL")
	cmd.PersistentFlags().Int("max-connections", 40, "Maximum allowed number of simultaneous HTTPS connections")

	return cmd
}

func runPolling(cmd *cobra.Command, args []string) error {
	// Get flags
	token, _ := cmd.Flags().GetString("token")
	envFile, _ := cmd.Flags().GetString("env")
	pollInterval, _ := cmd.Flags().GetInt("poll-interval")
	pollTimeout, _ := cmd.Flags().GetInt("poll-timeout")
	debug, _ := cmd.Flags().GetBool("debug")

	// Configure logging level
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// Load environment file if specified
	if envFile != "" {
		if err := godotenv.Load(envFile); err != nil {
			return fmt.Errorf("failed to load env file: %w", err)
		}
	}

	// Get token from environment if not provided
	if token == "" {
		token = os.Getenv("BOT_TOKEN")
		if token == "" {
			return fmt.Errorf("bot token is required (use --token or BOT_TOKEN env)")
		}
	}

	// Create bot
	b := bot.New(token, bot.WithLogger(log.Logger))

	// Set up signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start bot in polling mode
	log.Info().Str("mode", "polling").Msg("Starting bot")
	if err := b.StartPolling(
		time.Duration(pollInterval)*time.Second,
		time.Duration(pollTimeout)*time.Second,
	); err != nil {
		return fmt.Errorf("failed to start polling: %w", err)
	}

	// Wait for signal
	<-sigChan

	// Stop bot
	b.Stop()
	log.Info().Msg("Bot stopped")

	return nil
}

func runDev(cmd *cobra.Command, args []string) error {
	// Get flags
	token, _ := cmd.Flags().GetString("token")
	envFile, _ := cmd.Flags().GetString("env")
	port, _ := cmd.Flags().GetInt("port")
	webhookPath, _ := cmd.Flags().GetString("webhook-path")
	debug, _ := cmd.Flags().GetBool("debug")

	// Configure logging level
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// Load environment file if specified
	if envFile != "" {
		if err := godotenv.Load(envFile); err != nil {
			return fmt.Errorf("failed to load env file: %w", err)
		}
	}

	// Get token from environment if not provided
	if token == "" {
		token = os.Getenv("BOT_TOKEN")
		if token == "" {
			return fmt.Errorf("bot token is required (use --token or BOT_TOKEN env)")
		}
	}

	// Create bot
	b := bot.New(token, bot.WithLogger(log.Logger))

	// Set up signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start bot in webhook mode
	log.Info().Str("mode", "webhook").Int("port", port).Msg("Starting bot")
	if err := b.StartWebhook(webhookPath, port); err != nil {
		return fmt.Errorf("failed to start webhook: %w", err)
	}

	// Wait for signal
	<-sigChan

	// Stop bot
	b.Stop()
	log.Info().Msg("Bot stopped")

	return nil
}

func setWebhook(cmd *cobra.Command, args []string) error {
	// Get flags
	token, _ := cmd.Flags().GetString("token")
	url, _ := cmd.Flags().GetString("url")
	maxConnections, _ := cmd.Flags().GetInt("max-connections")

	// Get token from environment if not provided
	if token == "" {
		token = os.Getenv("BOT_TOKEN")
		if token == "" {
			return fmt.Errorf("bot token is required (use --token or BOT_TOKEN env)")
		}
	}

	// Get URL from environment if not provided
	if url == "" {
		url = os.Getenv("WEBHOOK_URL")
		if url == "" {
			return fmt.Errorf("webhook URL is required (use --url or WEBHOOK_URL env)")
		}
	}

	// Create bot
	b := bot.New(token, bot.WithLogger(log.Logger))

	// Set webhook
	log.Info().Str("url", url).Msg("Setting webhook")
	if err := b.SetWebhook(url, maxConnections, []string{"message", "callback_query", "inline_query"}); err != nil {
		return fmt.Errorf("failed to set webhook: %w", err)
	}

	// Get webhook info
	info, err := b.GetWebhookInfo()
	if err != nil {
		return fmt.Errorf("failed to get webhook info: %w", err)
	}

	log.Info().Interface("info", info).Msg("Webhook set successfully")
	return nil
}

func deleteWebhook(cmd *cobra.Command, args []string) error {
	// Get flags
	token, _ := cmd.Flags().GetString("token")

	// Get token from environment if not provided
	if token == "" {
		token = os.Getenv("BOT_TOKEN")
		if token == "" {
			return fmt.Errorf("bot token is required (use --token or BOT_TOKEN env)")
		}
	}

	// Create bot
	b := bot.New(token, bot.WithLogger(log.Logger))

	// Delete webhook
	log.Info().Msg("Deleting webhook")
	if err := b.DeleteWebhook(); err != nil {
		return fmt.Errorf("failed to delete webhook: %w", err)
	}

	log.Info().Msg("Webhook deleted successfully")
	return nil
}

func getWebhookInfo(cmd *cobra.Command, args []string) error {
	// Get flags
	token, _ := cmd.Flags().GetString("token")

	// Get token from environment if not provided
	if token == "" {
		token = os.Getenv("BOT_TOKEN")
		if token == "" {
			return fmt.Errorf("bot token is required (use --token or BOT_TOKEN env)")
		}
	}

	// Create bot
	b := bot.New(token, bot.WithLogger(log.Logger))

	// Get webhook info
	info, err := b.GetWebhookInfo()
	if err != nil {
		return fmt.Errorf("failed to get webhook info: %w", err)
	}

	log.Info().Interface("info", info).Msg("Webhook info")
	return nil
}
