package main

import (
	"fmt"

	"github.com/go-telegram/bot/models"
)

// Keyboards

// SaveToNotesButton returns a keyboard for saving a message to notes
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

// RequestLocation returns a keyboard for requesting location
func RequestLocation(messageID int) models.InlineKeyboardMarkup {
	btns := [][]models.InlineKeyboardButton{
		{
			{
				Text:         "📍 Share Location",
				CallbackData: fmt.Sprintf("share_location:%d", messageID),
			},
		},
		{
			{
				Text:         "✏️ Manual Input",
				CallbackData: "manual_location",
			},
		},
	}
	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: btns,
	}

	return *kb
}
