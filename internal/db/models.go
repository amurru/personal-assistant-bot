package db

import "time"

// User represents bot user information
type User struct {
	ID       int64     `json:"id"`                // Telegram User ID
	Name     string    `json:"name"`              // User Name
	City     string    `json:"city,omitempty"`    // User City
	Country  string    `json:"country,omitempty"` // User Country
	Phone    string    `json:"phone,omitempty"`   // User Phone Number
	Language string    `json:"language"`          // User Preferred Language
	Units    string    `json:"units"`             // User Preferred Units (i.e. metric)
	IsActive bool      `json:"is_active"`
	JoinedAt time.Time `json:"joined_at"`
}

// Note represents personal note information
type Note struct {
	ID        int       `json:"id,omitempty"`
	Text      string    `json:"text"`
	Owner     int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
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

// Quote represents inspirational quote information
type Quote struct {
	Text     string `json:"text"`
	Author   string `json:"author"`
	Source   string `json:"source"`
	URL      string `json:"url"`
	Language string `json:"lang"`
}

// LocationInfo represents location information
type LocationInfo struct {
	Country string `json:"country,omitempty"`
	City    string `json:"city,omitempty"`
	State   string `json:"state,omitempty"`
	Zip     string `json:"zip,omitempty"`
	Lat     string `json:"lat,omitempty"`
	Lon     string `json:"lon,omitempty"`
}

// UserStateInfo represents user state information
type UserStateInfo struct {
	ActiveCommand   string
	// can be previous message id, or something to edit ...etc
	CommandArgument any
}
