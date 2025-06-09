package main

import (
	"testing"
)

func TestGetWeatherInfo(t *testing.T) {
	result := GetWeatherInfo(
		"Jableh",
		"Syria",
		"metric",
	)
	if result == nil {
		t.Fatal("WeatherInfo is nil")
	}

	if result.City != "Jableh" {
		t.Errorf("Expected city Jableh, but got %s", result.City)
	}

	if result.Country != "Syria" {
		t.Errorf("Expected country Syria, but got %s", result.Country)
	}
}

func TestGetQuote(t *testing.T) {
	quote := GetQuote("en")
	if quote == nil {
		t.Fatal("Quote is nil")
	}

	if quote.Text == "" {
		t.Errorf("Quote text should not be empty")
	}

	if quote.Author == "" {
		t.Errorf("Quote author should not be empty")
	}
}
