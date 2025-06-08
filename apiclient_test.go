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
