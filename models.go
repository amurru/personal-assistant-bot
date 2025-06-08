package main

import "time"

// User represents bot user information
type User struct {
	ID       int64     `json:"id"`       // Telegram User ID
	Name     string    `json:"name"`     // User Name
	City     string    `json:"city"`     // User City
	Country  string    `json:"country"`  // User Country
	Phone    string    `json:"phone"`    // User Phone Number
	Language string    `json:"language"` // User Preferred Language
	Units    string    `json:"units"`    // User Preferred Units (i.e. metric)
	IsActive bool      `json:"is_active"`
	JoinedAt time.Time `json:"joined_at"`
}

// WeatherInfo represents weather forecast information
type WeatherInfo struct {
	Temp               string `json:"temp"`
	FeelsLike          string `json:"feels_like"`
	WeatherDescription string `json:"weather_description"`
	UVIndex            string `json:"uv_index"`
	Wind               string `json:"wind"`
	Precipitation      string `json:"precipitation"`
	Humidity           string `json:"humidity"`
	Pressure           string `json:"pressure"`
	Clouds             string `json:"clouds"`
	Visibility         string `json:"visibility"`
	City               string `json:"city"`
	Country            string `json:"country"`
	Units              string `json:"units"`
}
