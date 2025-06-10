package main

import (
	"context"
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
