package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Yashwanth-Kumar-26/teleriha/pkg/bot"
	"github.com/Yashwanth-Kumar-26/teleriha/pkg/gateway"
	"github.com/Yashwanth-Kumar-26/teleriha/pkg/store"
	"github.com/Yashwanth-Kumar-26/teleriha/plugins/broadcast"
	"github.com/Yashwanth-Kumar-26/teleriha/plugins/scenes"
	"github.com/Yashwanth-Kumar-26/teleriha/plugins/session"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	
	if err := godotenv.Load(".env"); err != nil {
		log.Warn().Msg("No .env file found, relying on environment variables")
	}

	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatal().Msg("BOT_TOKEN is required")
	}

	sqliteStore, err := store.NewSQLiteStore("bot_test.db")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open SQLite database")
	}
	defer func() { _ = sqliteStore.Close() }()

	app := bot.NewApplicationBuilder().
		Token(token).
		Opts(
			bot.WithLogger(log.Logger),
			bot.WithStore(sqliteStore),
		).
		Build()

	sceneManager := scenes.NewSceneManager()

	onboardingScene := scenes.NewScene("onboarding",
		func(ctx *bot.Context) error {
			_, err := ctx.Reply("Welcome to onboarding! What is your name?")
			return err
		},
		func(ctx *bot.Context) error {
			name := ctx.Message.Text
			sess := session.GetSession(ctx)
			if sess != nil {
				_ = sess.Set("name", []byte(name))
			}
			_, err := ctx.Reply(fmt.Sprintf("Nice to meet you, %s! Now, what is your email address?", name))
			return err
		},
		func(ctx *bot.Context) error {
			email := ctx.Message.Text
			sess := session.GetSession(ctx)
			var name string
			if sess != nil {
				nameBytes, _ := sess.Get("name")
				name = string(nameBytes)
			}
			_, err := ctx.Reply(fmt.Sprintf("Thank you, %s! Your registered email is %s. Onboarding complete!", name, email))
			if err != nil {
				return err
			}
			return sceneManager.Leave(ctx)
		},
	)
	sceneManager.AddScene(onboardingScene)

	app.Bot.Router.Use(session.Middleware(sqliteStore, 24*time.Hour))
	app.Bot.Router.Use(sceneManager.Middleware())
	app.Bot.Router.Use(bot.Recover(log.Logger))
	app.Bot.Router.Use(bot.Logger(log.Logger))

	broadcaster := broadcast.New(app.Bot)

	app.JobQueue.RunRepeating(func(ctx *bot.JobContext) error {
		log.Info().Msg("Periodic background job running...")
		return nil
	}, 30*time.Second, time.Now())

	app.AddHandler(bot.CommandHandler("/start", func(ctx *bot.Context) error {
		sess := session.GetSession(ctx)
		var visits int
		if sess != nil {
			val, _ := sess.Get("visits")
			if len(val) > 0 {
				visits, _ = strconv.Atoi(string(val))
			}
			visits++
			_ = sess.Set("visits", []byte(strconv.Itoa(visits)))
		}

		userData := ctx.UserData()
		if userData != nil {
			userData["last_command"] = "/start"
		}

		msg := fmt.Sprintf("Welcome to the Full-Stack Test Bot!\n\nYou have visited %d times.\n\nCommands:\n/paginate - Test Pagination Keyboard\n/invoice - Test Payments API (Stars)\n/broadcast - Test Broadcasting\n/onboard - Test FSM Scenes\n/gateway - Test Verification Gateway API\n\nYou can also send a Photo or Document to test logical message filters!", visits)
		_, err := ctx.Reply(msg)
		return err
	}))

	app.AddHandler(bot.CommandHandler("/paginate", func(ctx *bot.Context) error {
		pb := bot.NewPaginationBuilder(5, 1, "page_")
		keyboard := pb.Build()
		_, err := ctx.Reply("Pagination Example (Page 1):", bot.WithReplyMarkup(keyboard))
		return err
	}))

	app.AddHandler(bot.CallbackQueryHandler("page_", func(ctx *bot.Context) error {
		data := ctx.CallbackData()
		pageStr := data[len("page_"):]
		page, _ := strconv.Atoi(pageStr)

		pb := bot.NewPaginationBuilder(5, page, "page_")
		keyboard := pb.Build()

		_, err := ctx.EditMessageText(fmt.Sprintf("Pagination Example (Page %d):", page), bot.WithNewReplyMarkup(keyboard))
		_ = ctx.AnswerCallback("Navigated to page "+pageStr, false)
		return err
	}))

	app.AddHandler(bot.CommandHandler("/onboard", func(ctx *bot.Context) error {
		return sceneManager.Enter(ctx, "onboarding")
	}))

	app.AddHandler(bot.MessageHandler(bot.FilterPhoto(), func(ctx *bot.Context) error {
		_, err := ctx.Reply("Nice photo! The photo filter works.")
		return err
	}))

	app.AddHandler(bot.MessageHandler(bot.FilterDocument(), func(ctx *bot.Context) error {
		_, err := ctx.Reply("Received a document! The document filter works.")
		return err
	}))

	app.AddHandler(bot.CommandHandler("/invoice", func(ctx *bot.Context) error {
		prices := []bot.LabeledPrice{
			{Label: "Test Product", Amount: 50},
		}
		
		_, err := app.Bot.SendInvoice(
			ctx.ChatID(),
			"Full-Stack Bot T-Shirt",
			"A highly premium t-shirt for Telegram bot developers.",
			"payload_123",
			"",
			"XTR",
			prices,
		)
		if err != nil {
			_, _ = ctx.Reply(fmt.Sprintf("Error sending invoice: %v (Note: requires valid provider token in real usage)", err))
		}
		return nil
	}))

	app.AddHandler(bot.PreCheckoutQueryHandler(func(ctx *bot.Context) error {
		return app.Bot.AnswerPreCheckoutQuery(ctx.PreCheckoutQuery.ID, true, "")
	}))

	app.AddHandler(bot.CommandHandler("/broadcast", func(ctx *bot.Context) error {
		users := []int64{ctx.SenderID()}
		status := broadcaster.Broadcast(context.Background(), users, "This is a test broadcast message!")
		_, err := ctx.Reply(fmt.Sprintf("Broadcast initiated for %d users.", status.Total))
		return err
	}))

	app.AddHandler(bot.CommandHandler("/gateway", func(ctx *bot.Context) error {
		gatewayToken := os.Getenv("GATEWAY_TOKEN")
		if gatewayToken == "" {
			_, err := ctx.Reply("Please set GATEWAY_TOKEN in your .env to test the Telegram Gateway API.")
			return err
		}
		client := gateway.NewClient(gatewayToken)
		res, err := client.CheckSendAbility(context.Background(), "+1234567890")
		if err != nil {
			_, err = ctx.Reply(fmt.Sprintf("Gateway checkSendAbility failed: %v", err))
			return err
		}
		_, err = ctx.Reply(fmt.Sprintf("Gateway checkSendAbility response: ok=%t, request_id=%s", res.Ok, res.Result.RequestID))
		return err
	}))

	log.Info().Msg("Starting Full-Stack Test Bot...")
	if err := app.RunPolling(); err != nil {
		log.Fatal().Err(err).Msg("Failed to start application")
	}

	select {}
}
