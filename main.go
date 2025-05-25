package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	_ "github.com/dotenv-org/godotenvvault/autoload"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

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

	opts := []bot.Option{
		bot.WithDebug(),
		bot.WithMessageTextHandler("/start", bot.MatchTypeExact, startHandler),
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
			break
		}
	}

	b.Start(ctx)
}

func startHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Hello World!",
	})
}
