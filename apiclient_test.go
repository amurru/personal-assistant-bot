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
		t.Error("WeatherInfo is nil")
	}
}
