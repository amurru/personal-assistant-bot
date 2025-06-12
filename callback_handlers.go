package main

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/amurru/personal-assistant-bot/internal/db"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func saveToNotesHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})
	owner, err := pers.GetUser(update.CallbackQuery.From.ID)
	note := db.Note{
		Text:      update.CallbackQuery.Message.Message.Text,
		Owner:     owner.ID,
		CreatedAt: time.Now(),
	}
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.CallbackQuery.Message.Message.Chat.ID,
			Text:   "Error occured. Contact developer!",
		})
		return
	}
	pers.AddNote(note)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.CallbackQuery.Message.Message.Chat.ID,
		Text:   "Saved! Check with /notes",
	})
}

func shareLocationHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.CallbackQuery.Message.Message.Chat.ID,
		Text:   "Please send me your location from the pin menu",
	})
	// Extract previous message id from share_location:ID to use as reference
	previousMessageID, err := strconv.Atoi(
		strings.Split(update.CallbackQuery.Data, ":")[1], // share_location:msgid
	)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.CallbackQuery.Message.Message.Chat.ID,
			Text:   "Error occured. Contact developer!",
		})
		return
	}
	// Create user state record if not exists
	if _, ok := userStates[update.CallbackQuery.From.ID]; !ok {
		userStates[update.CallbackQuery.From.ID] = &db.UserStateInfo{}
	}
	userStates[update.CallbackQuery.From.ID].PreviousMessageID = previousMessageID
	userStates[update.CallbackQuery.From.ID].ActiveCommand = "waiting_for_location"
}
