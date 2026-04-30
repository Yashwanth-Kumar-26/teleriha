// Package main provides an ecommerce bot example using TeleRiHa.
package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/Yashwanth-Kumar-26/teleriha/pkg/bot"
)

// Product represents a product in the catalog.
type Product struct {
	ID          string
	Name        string
	Description string
	Price       float64
	ImageURL    string
	Category    string
}

// Cart represents a user's shopping cart.
type Cart struct {
	Items map[string]int // product ID -> quantity
}

// Catalog of products
var catalog = map[string]Product{
	"p001": {
		ID:          "p001",
		Name:        "Wireless Headphones",
		Description: "High-quality wireless headphones with noise cancellation",
		Price:       199.99,
		ImageURL:    "https://example.com/images/headphones.jpg",
		Category:    "Electronics",
	},
	"p002": {
		ID:          "p002",
		Name:        "Smartphone",
		Description: "Latest smartphone with advanced features",
		Price:       799.99,
		ImageURL:    "https://example.com/images/smartphone.jpg",
		Category:    "Electronics",
	},
	"p003": {
		ID:          "p003",
		Name:        "Coffee Mug",
		Description: "Keep your drinks hot with this insulated mug",
		Price:       19.99,
		ImageURL:    "https://example.com/images/mug.jpg",
		Category:    "Home",
	},
	"p004": {
		ID:          "p004",
		Name:        "T-Shirt",
		Description: "Comfortable cotton t-shirt",
		Price:       24.99,
		ImageURL:    "https://example.com/images/tshirt.jpg",
		Category:    "Clothing",
	},
}

// In-memory cart storage
var carts = make(map[int64]*Cart)

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
	b.Router.On("/products", handleProducts)
	b.Router.On("/categories", handleCategories)
	b.Router.On("/category", handleCategory)
	b.Router.On("/add", handleAddToCart)
	b.Router.On("/cart", handleViewCart)
	b.Router.On("/checkout", handleCheckout)
	b.Router.On("/clear", handleClearCart)

	// Register callback query handlers
	b.Router.OnCallback("product_", handleProductCallback)
	b.Router.OnCallback("category_", handleCategoryCallback)
	b.Router.OnCallback("cart_", handleCartCallback)

	// Register a default handler
	b.Router.Default(handleDefault)

	// Start the bot in polling mode
	log.Info().Msg("Starting ecommerce bot in polling mode...")
	if err := b.StartPolling(1, 30); err != nil {
		log.Fatal().Err(err).Msg("Failed to start bot")
	}

	// Wait forever
	select {}
}

func handleStart(ctx *bot.Context) error {
	message := `🏪 *Welcome to TeleShop!* 🏪

Your personal shopping assistant on Telegram. Browse our catalog, add items to your cart, and checkout with ease.

Use /help to see all commands.`

	kb := bot.NewInlineKeyboardBuilder().
		AddRow(
			bot.InlineButton("🛍️ Browse Products", "command:products"),
			bot.InlineButton("📁 Categories", "command:categories"),
		).
		AddRow(
			bot.InlineButton("🛒 My Cart", "command:cart"),
			bot.InlineButton("❓ Help", "command:help"),
		).
		Build()

	_, err := ctx.Reply(message, bot.WithParseMode("Markdown"), bot.WithReplyMarkup(kb))
	return err
}

func handleHelp(ctx *bot.Context) error {
	helpText := `📋 *Available Commands:*

🛍️ *Shopping:*
/products - Browse all products
/categories - View product categories
/category <name> - View products in a category
/add <product_id> [quantity] - Add product to cart
/cart - View your shopping cart
/checkout - Checkout and place order
/clear - Clear your shopping cart

💡 *Tips:*
- Use buttons for quick navigation
- Add items to cart and checkout when ready
- Your cart is saved for this session`

	_, err := ctx.Reply(helpText, bot.WithParseMode("Markdown"))
	return err
}

func handleProducts(ctx *bot.Context) error {
	if len(catalog) == 0 {
		_, err := ctx.Reply("No products available.")
		return err
	}

	message := "*📦 Available Products:*\n\n"
	count := 0
	for id, product := range catalog {
		count++
		message += fmt.Sprintf("🏷️ *%s*\n%s - $%.2f\n[Add to Cart](/add_%s)\n\n",
			product.Name, product.Description, product.Price, id)
		if count >= 5 {
			message += "_and more... use /category to filter_\n"
			break
		}
	}

	kb := bot.NewInlineKeyboardBuilder().
		AddRow(bot.InlineButton("⬅️ Back to Main", "command:start")).
		Build()

	_, err := ctx.Reply(message, bot.WithParseMode("Markdown"), bot.WithReplyMarkup(kb))
	return err
}

func handleCategories(ctx *bot.Context) error {
	categories := make(map[string]bool)
	for _, product := range catalog {
		categories[product.Category] = true
	}

	if len(categories) == 0 {
		_, err := ctx.Reply("No categories available.")
		return err
	}

	message := "*📁 Product Categories:*\n"

	kb := bot.NewInlineKeyboardBuilder()
	for category := range categories {
		kb.AddRow(bot.InlineButton(category, "category:"+category))
	}
	kb.AddRow(bot.InlineButton("⬅️ Back", "command:start"))

	_, err := ctx.Reply(message, bot.WithParseMode("Markdown"), bot.WithReplyMarkup(kb.Build()))
	return err
}

func handleCategory(ctx *bot.Context) error {
	args := ctx.CommandArgs()
	if args == "" {
		_, err := ctx.Reply("Please specify a category. Usage: /category <name>")
		return err
	}

	category := strings.Title(strings.ToLower(args))

	// Find products in this category
	var products []Product
	for _, product := range catalog {
		if strings.EqualFold(product.Category, category) {
			products = append(products, product)
		}
	}

	if len(products) == 0 {
		_, err := ctx.Reply(fmt.Sprintf("No products found in category '%s'.", category))
		return err
	}

	message := fmt.Sprintf("*📦 Products in %s:*\n\n", category)
	for _, product := range products {
		message += fmt.Sprintf("🏷️ *%s*\n%.2f USD\n[Add to Cart](command:add_%s)\n\n",
			product.Name, product.Price, product.ID)
	}

	kb := bot.NewInlineKeyboardBuilder().
		AddRow(bot.InlineButton("⬅️ Back to Categories", "command:categories")).
		Build()

	_, err := ctx.Reply(message, bot.WithParseMode("Markdown"), bot.WithReplyMarkup(kb))
	return err
}

func handleProductCallback(ctx *bot.Context) error {
	callbackData := ctx.CallbackData()
	if !strings.HasPrefix(callbackData, "product:") {
		return nil
	}

	productID := strings.TrimPrefix(callbackData, "product:")
	product, ok := catalog[productID]
	if !ok {
		return ctx.AnswerCallback("Product not found", true)
	}

	message := fmt.Sprintf("*%s*\n\n%s\n\nPrice: *$%.2f*\n\nSelect an action:",
		product.Name, product.Description, product.Price)

	kb := bot.NewInlineKeyboardBuilder().
		AddRow(
			bot.InlineButton("➕ Add to Cart", "add:"+product.ID+":1"),
		).
		AddRow(
			bot.InlineButton("🔙 Back", "command:products"),
		).
		Build()

	// Edit the original message
	_, err := ctx.EditMessageText(message,
		bot.WithParseMode("Markdown"),
		bot.WithNewReplyMarkup(kb))
	if err != nil {
		// If we can't edit, just answer the callback
		return ctx.AnswerCallback("Viewing product details", false)
	}

	return ctx.AnswerCallback("", false)
}

func handleCategoryCallback(ctx *bot.Context) error {
	callbackData := ctx.CallbackData()
	if !strings.HasPrefix(callbackData, "category:") {
		return nil
	}

	category := strings.TrimPrefix(callbackData, "category:")

	// Find products in this category
	var products []Product
	for _, product := range catalog {
		if strings.EqualFold(product.Category, category) {
			products = append(products, product)
		}
	}

	if len(products) == 0 {
		return ctx.AnswerCallback("No products in this category", true)
	}

	message := fmt.Sprintf("*📦 %s:*\n\n", category)
	for _, product := range products {
		message += fmt.Sprintf("🏷️ *%s* - $%.2f\n[View](product:%s)\n\n",
			product.Name, product.Price, product.ID)
	}

	kb := bot.NewInlineKeyboardBuilder().
		AddRow(bot.InlineButton("⬅️ Back", "command:categories")).
		Build()

	// Edit the original message
	_, err := ctx.EditMessageText(message,
		bot.WithParseMode("Markdown"),
		bot.WithNewReplyMarkup(kb))
	if err != nil {
		return ctx.AnswerCallback("Showing category products", false)
	}

	return ctx.AnswerCallback("", false)
}

func handleAddToCart(ctx *bot.Context) error {
	args := ctx.CommandArgs()
	if args == "" {
		_, err := ctx.Reply("Please specify a product ID. Usage: /add <product_id> [quantity]")
		return err
	}

	parts := strings.Fields(args)
	productID := parts[0]
	quantity := 1

	if len(parts) > 1 {
		var err error
		quantity, err = strconv.Atoi(parts[1])
		if err != nil || quantity < 1 {
			_, err := ctx.Reply("Invalid quantity. Usage: /add <product_id> [quantity]")
			return err
		}
	}

	product, ok := catalog[productID]
	if !ok {
		_, err := ctx.Reply("Product not found. Use /products to browse available products.")
		return err
	}

	// Get or create cart for user
	userID := ctx.SenderID()
	if carts[userID] == nil {
		carts[userID] = &Cart{Items: make(map[string]int)}
	}

	// Add product to cart
	carts[userID].Items[productID] += quantity

	_, err := ctx.Reply(fmt.Sprintf("✅ Added %d x *%s* to your cart.\nTotal: $%.2f",
		quantity, product.Name, product.Price*float64(quantity)),
		bot.WithParseMode("Markdown"))
	return err
}

func handleViewCart(ctx *bot.Context) error {
	userID := ctx.SenderID()
	cart := carts[userID]

	if cart == nil || len(cart.Items) == 0 {
		_, err := ctx.Reply("🛒 Your cart is empty. Use /products to add items.")
		return err
	}

	message := "*🛒 Your Shopping Cart:*\n\n"
	total := 0.0
	count := 0

	for productID, quantity := range cart.Items {
		product := catalog[productID]
		if product.ID == "" {
			continue
		}
		itemTotal := product.Price * float64(quantity)
		total += itemTotal
		count += quantity
		message += fmt.Sprintf("%d x *%s* - $%.2f each = *$%.2f*\n",
			quantity, product.Name, product.Price, itemTotal)
	}

	message += fmt.Sprintf("\n📊 *Summary:* %d items, *Total: $%.2f*\n\n", count, total)

	kb := bot.NewInlineKeyboardBuilder().
		AddRow(
			bot.InlineButton("💳 Checkout", "checkout:confirm"),
			bot.InlineButton("🗑️ Clear Cart", "clear:confirm"),
		).
		AddRow(
			bot.InlineButton("⬅️ Continue Shopping", "command:products"),
		).
		Build()

	_, err := ctx.Reply(message, bot.WithParseMode("Markdown"), bot.WithReplyMarkup(kb))
	return err
}

func handleCheckout(ctx *bot.Context) error {
	userID := ctx.SenderID()
	cart := carts[userID]

	if cart == nil || len(cart.Items) == 0 {
		_, err := ctx.Reply("🛒 Your cart is empty. Add items before checking out.")
		return err
	}

	total := 0.0
	for productID, quantity := range cart.Items {
		product := catalog[productID]
		if product.ID == "" {
			continue
		}
		total += product.Price * float64(quantity)
	}

	message := fmt.Sprintf("*💳 Checkout Confirmation*\n\nOrder Total: *$%.2f*\n\nConfirm your order?", total)

	kb := bot.NewInlineKeyboardBuilder().
		AddRow(
			bot.InlineButton("✅ Confirm", "checkout:place_order"),
			bot.InlineButton("❌ Cancel", "checkout:cancel"),
		).
		Build()

	_, err := ctx.Reply(message, bot.WithParseMode("Markdown"), bot.WithReplyMarkup(kb))
	return err
}

func handleClearCart(ctx *bot.Context) error {
	userID := ctx.SenderID()

	if carts[userID] == nil || len(carts[userID].Items) == 0 {
		_, err := ctx.Reply("🛒 Your cart is already empty.")
		return err
	}

	message := "Are you sure you want to clear your cart?"

	kb := bot.NewInlineKeyboardBuilder().
		AddRow(
			bot.InlineButton("✅ Yes, Clear", "clear:confirmed"),
			bot.InlineButton("❌ Cancel", "clear:cancel"),
		).
		Build()

	_, err := ctx.Reply(message, bot.WithReplyMarkup(kb))
	return err
}

func handleCartCallback(ctx *bot.Context) error {
	callbackData := ctx.CallbackData()

	// Handle checkout callbacks
	if strings.HasPrefix(callbackData, "checkout:") {
		action := strings.TrimPrefix(callbackData, "checkout:")
		userID := ctx.SenderID()
		cart := carts[userID]

		switch action {
		case "place_order":
			if cart != nil && len(cart.Items) > 0 {
				// Process order
				total := 0.0
				for productID, quantity := range cart.Items {
					product := catalog[productID]
					if product.ID != "" {
						total += product.Price * float64(quantity)
					}
				}

				// Clear cart
				carts[userID] = &Cart{Items: make(map[string]int)}

				return ctx.AnswerCallback(fmt.Sprintf("✅ Order placed! Total: $%.2f", total), true)
			}
			return ctx.AnswerCallback("❌ Cart is empty", true)

		case "confirm":
			return handleCheckout(ctx)

		case "cancel":
			return ctx.AnswerCallback("Checkout cancelled", true)
		}
	}

	// Handle clear callbacks
	if strings.HasPrefix(callbackData, "clear:") {
		action := strings.TrimPrefix(callbackData, "clear:")
		userID := ctx.SenderID()

		switch action {
		case "confirmed":
			carts[userID] = &Cart{Items: make(map[string]int)}
			return ctx.AnswerCallback("✅ Cart cleared", true)

		case "cancel":
			return ctx.AnswerCallback("Cart not cleared", true)
		}
	}

	return nil
}

func handleDefault(ctx *bot.Context) error {
	text := ctx.Text()
	if text == "" {
		return nil
	}

	// Handle command-like messages without the slash
	lowerText := strings.ToLower(text)
	if strings.HasPrefix(lowerText, "start") || strings.HasPrefix(lowerText, "help") {
		return handleHelp(ctx)
	}

	_, err := ctx.Reply(fmt.Sprintf("I don't understand '%s'. Use /help to see available commands.", text))
	return err
}
