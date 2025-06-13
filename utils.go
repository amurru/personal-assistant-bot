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
				CallbackData: fmt.Sprintf("manual_location:%d", messageID),
			},
		},
	}
	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: btns,
	}

	return *kb
}

func NotesActionButtons() models.InlineKeyboardMarkup {
	btns := [][]models.InlineKeyboardButton{
		{
			{
				Text:         "📋 Add",
				CallbackData: "notes_add",
			},
			{
				Text:         "📝 Edit",
				CallbackData: "notes_edit",
			},
			{
				Text:         "🗑️ Delete",
				CallbackData: "notes_delete",
			},
		},
	}
	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: btns,
	}

	return *kb
}

func ProfileActionButtons() models.InlineKeyboardMarkup {
	btns := [][]models.InlineKeyboardButton{
		{
			{
				Text:         "👤 Edit Name",
				CallbackData: "profile_change_name",
			},
			{
				Text:         "📞 Edit Phone",
				CallbackData: "profile_change_phone",
			},
			{
				Text:         "🌐 Edit Language",
				CallbackData: "profile_change_language",
			},
		},
		{
			{
				Text:         "📍 Edit Address",
				CallbackData: "profile_change_address",
			},
			{
				Text:         "📏 Toggle Units",
				CallbackData: "profile_change_units",
			},
		},
	}
	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: btns,
	}
	return *kb
}

