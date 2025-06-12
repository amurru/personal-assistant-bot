// Package main is the main package for the personal assistant bot
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

var (
	pers       db.Persistence
	userStates = make(map[int64]*db.UserStateInfo)
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
		// call-back handlers
		bot.WithCallbackQueryDataHandler("save_to_notes", bot.MatchTypeExact, saveToNotesHandler),
		bot.WithCallbackQueryDataHandler(
			"share_location:",
			bot.MatchTypePrefix,
			shareLocationHandler,
		),
		bot.WithCallbackQueryDataHandler("manual_location", bot.MatchTypeExact, locationHandler),

		// default handler
		bot.WithDefaultHandler(defaultHandler),
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

		// Onboarding Experience
		// Register user
		user := db.NewUserObject()
		user.ID = update.Message.From.ID
		user.Name = fmt.Sprintf(
			"%s %s",
			update.Message.From.FirstName,
			update.Message.From.LastName,
		)
		user.Language = update.Message.From.LanguageCode
		user.JoinedAt = time.Now()
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   fmt.Sprintf("Hello %s, I'm here to help you!", update.Message.From.FirstName),
		})

		// 1. Survey user information and explain how they are used
		// i.e city, country, for weather (optional), units (default is metric = m)
		locationRequestText := "To provide you with weather updates, can I use your location?"
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Message.Chat.ID,
			Text:        locationRequestText,
			ReplyMarkup: RequestLocation(update.Message.ID),
		})

		// 2. Ask user to confirm

		// 3. Register user in DB
		err := pers.AddUser(user)
		if err != nil {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Error occured. Contact developer!",
			})
			return
		}
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

func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	if _, ok := userStates[update.Message.From.ID]; !ok {
		// ignore for now
		return
	}
	switch userStates[update.Message.From.ID].ActiveCommand {
	case "waiting_for_location":
		if update.Message.Location != nil {
			locationInfo, err := GetLocationInformation(
				update.Message.Location.Latitude,
				update.Message.Location.Longitude,
			)
			if err != nil {
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: update.Message.Chat.ID,
					Text:   "Error getting location information. Please try again later.",
				})
				return
			}

			// update user location
			user, err := pers.GetUser(update.Message.From.ID)
			if err != nil {
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: update.Message.Chat.ID,
					Text:   "Error getting user information. Please try again later.",
				})
				return
			}
			user.City = locationInfo.City
			user.Country = locationInfo.Country
			err = pers.UpdateUser(user)
			if err != nil {
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: update.Message.Chat.ID,
					Text:   "Error updating user information. Please try again later.",
				})
				return
			}

			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Your location has been updated. Please confirm your details.",
			})

			// FIX: message has not been deleted
			// delete message requesting location
			b.DeleteMessage(ctx, &bot.DeleteMessageParams{
				ChatID:    update.Message.Chat.ID,
				MessageID: userStates[user.ID].PreviousMessageID,
			})

			// remove user's state
			delete(userStates, user.ID)
			return

		}
	}
}

func locationHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})

	user, err := pers.GetUser(update.CallbackQuery.From.ID)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.CallbackQuery.Message.Message.Chat.ID,
			Text:   "An error occurred. Please try again later.",
		})
		return
	}
	if update.Message.Location != nil {
		// User shared location
		// TODO: we will call a helper to get city and country from geolocation data
		// update.Message.Location.Latitude
		// update.Message.Location.Longitude

		_ = user
		GetLocationInformation(update.Message.Location.Latitude, update.Message.Location.Longitude)
		// Ask for preferred units
		unitsKeyboard := &models.ReplyKeyboardMarkup{
			Keyboard: [][]models.KeyboardButton{
				{
					{
						Text: "Metric (Celsius)",
					},
				},
				{
					{
						Text: "Imperial (Fahrenheit)",
					},
				},
			},
			ResizeKeyboard:  true,
			OneTimeKeyboard: true,
		}

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
			Text:        "Great! Now, do you prefer metric (Celsius) or imperial (Fahrenheit) units?",
			ReplyMarkup: unitsKeyboard,
		})

		userStates[user.ID].ActiveCommand = "waiting_for_units"

	} else if update.Message.Text == "No, enter manually" {
		// User chose to enter location manually
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.CallbackQuery.Message.Message.Chat.ID,
			Text:   "Please enter your city:",
		})

		userStates[user.ID].ActiveCommand = "waiting_for_city"
	} else {
		// Handle city and units input
		state, ok := userStates[user.ID]
		if !ok {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.CallbackQuery.Message.Message.Chat.ID,
				Text:   "Sorry, I didn't understand that. Please start again with /start",
			})
			return
		}

		switch state.ActiveCommand {
		case "waiting_for_city":
			user.City = update.Message.Text
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.CallbackQuery.Message.Message.Chat.ID,
				Text:   "Please enter your country:",
			})
			userStates[user.ID].ActiveCommand = "waiting_for_country"
		case "waiting_for_country":
			user.Country = update.Message.Text

			// Ask for preferred units
			unitsKeyboard := &models.ReplyKeyboardMarkup{
				Keyboard: [][]models.KeyboardButton{
					{
						{
							Text: "Metric (Celsius)",
						},
					},
					{
						{
							Text: "Imperial (Fahrenheit)",
						},
					},
				},
				ResizeKeyboard:  true,
				OneTimeKeyboard: true,
			}

			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
				Text:        "Great! Now, do you prefer metric (Celsius) or imperial (Fahrenheit) units?",
				ReplyMarkup: unitsKeyboard,
			})
			userStates[user.ID].ActiveCommand = "waiting_for_units"
		case "waiting_for_units":
			if update.Message.Text == "Metric (Celsius)" {
				user.Units = "metric"
			} else if update.Message.Text == "Imperial (Fahrenheit)" {
				user.Units = "imperial"
			} else {
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: update.CallbackQuery.Message.Message.Chat.ID,
					Text:   "Invalid input. Please select from the options provided.",
				})
				return
			}

			// Confirm user info
			confirmationText := fmt.Sprintf(
				"Please confirm your details:\nCity: %s\nCountry: %s\nUnits: %s",
				user.City,
				user.Country,
				user.Units,
			)

			confirmKeyboard := &models.ReplyKeyboardMarkup{
				Keyboard: [][]models.KeyboardButton{
					{
						{
							Text: "Confirm",
						},
					},
					{
						{
							Text: "Cancel",
						},
					},
				},
				ResizeKeyboard:  true,
				OneTimeKeyboard: true,
			}

			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
				Text:        confirmationText,
				ReplyMarkup: confirmKeyboard,
			})
			userStates[user.ID].ActiveCommand = "waiting_for_confirmation"
		case "waiting_for_confirmation":
			if update.Message.Text == "Confirm" {
				err := pers.UpdateUser(user)
				if err != nil {
					b.SendMessage(ctx, &bot.SendMessageParams{
						ChatID: update.CallbackQuery.Message.Message.Chat.ID,
						Text:   "Error updating user information. Please contact the developer.",
					})
					return
				}
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: update.CallbackQuery.Message.Message.Chat.ID,
					Text:   "Your information has been saved. Welcome!",
				})
			} else if update.Message.Text == "Cancel" {
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: update.CallbackQuery.Message.Message.Chat.ID,
					Text:   "Onboarding cancelled. You can start again with /start",
				})
			}
			delete(userStates, user.ID)
		}
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

	// Get user location and preferences from database
	user, err := pers.GetUser(update.Message.From.ID)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Error occured. Contact developer!",
		})
		return
	}
	w := GetWeatherInfo(user.City, user.Country, user.Units)

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
	notes, err := pers.GetUserNotes(update.Message.From.ID)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Error occured. Contact developer!",
		})
		return
	}
	var notesText string
	for i, note := range notes {
		notesText += fmt.Sprintf("%d. %s\n\n", i+1, note.Text)
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Your Notes:\n-------------\n\n" + notesText,
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
		ChatID:      update.Message.Chat.ID,
		Text:        quoteText,
		ParseMode:   models.ParseModeMarkdownV1,
		ReplyMarkup: SaveToNotesButton(),
	})
}
