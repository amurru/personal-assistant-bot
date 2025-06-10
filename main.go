package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/amurru/personal-assistant-bot/internal/db"
	_ "github.com/dotenv-org/godotenvvault/autoload"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

var pers db.Persistence

func main() {
	tgToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if tgToken == "" {
		fmt.Println("TELEGRAM_BOT_TOKEN is not set")
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
	)
	defer cancel()

	/*
		Commands:
		start - Start Interacting with Bot
		brief - Show a summary of upcoming activities
		remind - Add New Reminder
		weather - Query Weather Forecast
		calendar - Manage Calendar
		notes - Manage Personal Notes
		inspire - Get Inspirational Quote
		request - Request New Features
		help - Show Help Info
	*/
	opts := []bot.Option{
		bot.WithMessageTextHandler("/start", bot.MatchTypeExact, startHandler),
		bot.WithMessageTextHandler("/help", bot.MatchTypeExact, helpHandler),
		bot.WithMessageTextHandler("/calendar", bot.MatchTypeExact, calendarHandler),
		bot.WithMessageTextHandler("/remind", bot.MatchTypeExact, remindHandler),
		bot.WithMessageTextHandler("/request", bot.MatchTypeExact, requestHandler),
		bot.WithMessageTextHandler("/weather", bot.MatchTypeExact, weatherHandler),
		bot.WithMessageTextHandler("/notes", bot.MatchTypeExact, notesHandler),
		bot.WithMessageTextHandler("/brief", bot.MatchTypeExact, briefHandler),
		bot.WithMessageTextHandler("/inspire", bot.MatchTypeExact, inspireHandler),
	}

	// check if debug mode
	debug := os.Getenv("BOT_DEBUG")
	if debug == "true" {
		opts = append(opts, bot.WithDebug())
	}

	// custom bot API server
	serverURL := os.Getenv("BOT_API_SERVER")
	if serverURL != "" {
		opts = append(opts, bot.WithServerURL(serverURL))
	}

	var b *bot.Bot
	var err error
	maxTries := 5
	for {
		b, err = bot.New(tgToken, opts...)
		if err != nil {
			log.Printf("Error Launching Bot: %v", err)
			if maxTries > 0 {
				log.Println("Retrying...")
				maxTries--
			} else {
				log.Fatal("Could not connect to Telegram")
			}
		} else {
			log.Println("Bot Launched")
			break
		}
	}

	// prepare database
	pers = db.InstanceOrNew()

	b.Start(ctx)
}

func startHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	// Check if user is already registered
	if !pers.IsKnownUser(update.Message.From.ID) {
		// Register user
		// 1. Survey user information and explain how they are used
		// i.e city, country, for weather (optional). units (default is metric = m)
		// 2. Ask user to confirm
		// 3. Register user in DB
		user := db.NewUserObject()
		user.ID = update.Message.From.ID
		user.Name = fmt.Sprintf(
			"%s %s",
			update.Message.From.FirstName,
			update.Message.From.LastName,
		)
		user.Country = "Spain"
		user.City = "Madrid"
		user.Language = "es"
		user.Units = "m"
		user.JoinedAt = time.Now()
		err := pers.AddUser(user)
		if err != nil {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Error occured. Contact developer!",
			})
			return
		}
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   fmt.Sprintf("Hello %s, I'm here to help you!", update.Message.From.FirstName),
		})
	} else {
		user, err := pers.GetUser(update.Message.From.ID)
		if err != nil {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "An Error occurred. Try again later!\nIf persisted please contact developer",
			})
			return
		}
		welcomeMessage := fmt.Sprintf("Hello %s, welcome back!",
			strings.Split(user.Name, " ")[0],
		)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   welcomeMessage,
		})
	}
}

func helpHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Help Text",
	})
}
func calendarHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Calendar",
	})
}
func remindHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Reminder",
	})
}

func requestHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Request",
	})
}

func weatherHandler(ctx context.Context, b *bot.Bot, update *models.Update) {

	// get user location and preferences from database
	// but for now, just use New York, US, lang, imperial
	w := GetWeatherInfo("New York", "US", "imperial")

	weathernReport := fmt.Sprintf(
		"Temperature: %s\nFeels Like: %s\nUV Index: %s\nWind: %s\nPrecipitation: %s\nHumidity: %s\nPressure: %s\nClouds: %s\nVisibility: %s\nCity: %s\nCountry: %s\nUnits: %s",
		w.Temp,
		w.FeelsLike,
		w.UVIndex,
		w.Wind,
		w.Precipitation,
		w.Humidity,
		w.Pressure,
		w.Clouds,
		w.Visibility,
		w.City,
		w.Country,
		w.Units,
	)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   weathernReport,
	})
}

func notesHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Notes",
	})
}

func briefHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Brief",
	})
}

func inspireHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	quote := GetQuote("en")
	if quote == nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Error Getting Quote",
		})
		return
	}
	quoteText := fmt.Sprintf("*🙶%s🙸*\n—%s", quote.Text, quote.Author)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      quoteText,
		ParseMode: models.ParseModeMarkdownV1,
	})
}
