package main

import (
	"github.com/go-telegram/bot/models"
)

// Keyboards
func SaveToNotesButton() models.InlineKeyboardMarkup {
	btns := [][]models.InlineKeyboardButton{
		{
			{
				Text:         "💾 Save to Notes",
				CallbackData: "save_to_notes",
			},
		},
	}
	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: btns,
	}

	return *kb
}
