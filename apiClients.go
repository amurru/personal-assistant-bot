package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/amurru/personal-assistant-bot/internal/db"
)

// GetWeatherInfo returns a pointer to a WeatherInfo struct
// with weather information for the given city and country.
func GetWeatherInfo(city, country, units string) *db.WeatherInfo {
	url := fmt.Sprintf(
		"https://wttr.in/%s-%s?format=j1",
		country,
		city,
	)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil
	}
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{
		Timeout: 30 * time.Second,
	}
	res, err := client.Do(req)
	if err != nil {
		return nil
	}
	// Unmarshal response body to anonymous struct
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil
	}
	// Define in-line struct for unmarshalling
	var weatherData struct {
		CurrentCondition []struct {
			FeelsLikeC      string `json:"FeelsLikeC"`
			FeelsLikeF      string `json:"FeelsLikeF"`
			TempC           string `json:"temp_C"`
			TempF           string `json:"temp_F"`
			UvIndex         string `json:"uvIndex"`
			WindDir16Point  string `json:"winddir16Point"`
			WindDirDegree   string `json:"winddirDegree"`
			WindspeedKmph   string `json:"windspeedKmph"`
			WindspeedMiles  string `json:"windspeedMiles"`
			PrecipInches    string `json:"precipInches"`
			PrecipMM        string `json:"precipMM"`
			Humidity        string `json:"humidity"`
			Pressure        string `json:"pressure"`
			PressureInches  string `json:"pressureInches"`
			Clouds          string `json:"cloudcover"`
			Visibility      string `json:"visibility"`
			VisibilityMiles string `json:"visibilityMiles"`
			ObservationTime string `json:"observation_time"`
		} `json:"current_condition"`
	}
	json.Unmarshal(resBody, &weatherData)
	// Populate weatherInfo struct based on units
	var weatherInfo db.WeatherInfo
	if len(weatherData.CurrentCondition) == 0 {
		log.Println("No weather data found for city")
		return nil
	}
	if units == "imperial" || units == "i" {
		weatherInfo.Temp = weatherData.CurrentCondition[0].TempF
		weatherInfo.FeelsLike = weatherData.CurrentCondition[0].FeelsLikeF
		weatherInfo.UVIndex = weatherData.CurrentCondition[0].UvIndex
		weatherInfo.Wind = fmt.Sprintf(
			"%s %s",
			weatherData.CurrentCondition[0].WindspeedMiles,
			weatherData.CurrentCondition[0].WindDir16Point,
		)
		weatherInfo.Precipitation = weatherData.CurrentCondition[0].PrecipInches
		weatherInfo.Humidity = weatherData.CurrentCondition[0].Humidity
		weatherInfo.Pressure = weatherData.CurrentCondition[0].PressureInches
		weatherInfo.Clouds = weatherData.CurrentCondition[0].Clouds
		weatherInfo.Visibility = weatherData.CurrentCondition[0].VisibilityMiles
	} else if units == "metric" || units == "m" {
		weatherInfo.Temp = weatherData.CurrentCondition[0].TempC
		weatherInfo.FeelsLike = weatherData.CurrentCondition[0].FeelsLikeC
		weatherInfo.UVIndex = weatherData.CurrentCondition[0].UvIndex
		weatherInfo.Wind = fmt.Sprintf(
			"%s %s",
			weatherData.CurrentCondition[0].WindspeedKmph,
			weatherData.CurrentCondition[0].WindDirDegree,
		)
		weatherInfo.Precipitation = weatherData.CurrentCondition[0].PrecipMM
		weatherInfo.Humidity = weatherData.CurrentCondition[0].Humidity
		weatherInfo.Pressure = weatherData.CurrentCondition[0].Pressure
		weatherInfo.Clouds = weatherData.CurrentCondition[0].Clouds
		weatherInfo.Visibility = weatherData.CurrentCondition[0].Visibility
	}
	weatherInfo.City = city
	weatherInfo.Country = country
	weatherInfo.Units = units

	return &weatherInfo
}

// GetQuote gives a random inspirational quote in the given language
func GetQuote(lang string) *db.Quote {
	url := "https://thequoteshub.com/api/"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil
	}
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{
		Timeout: 30 * time.Second,
	}
	res, err := client.Do(req)
	if err != nil {
		return nil
	}
	// Unmarshal response body to anonymous struct
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil
	}
	var quoteData struct {
		Text   string `json:"text"`
		Author string `json:"author"`
	}

	json.Unmarshal(resBody, &quoteData)

	// For now, just return English quote
	_ = lang
	return &db.Quote{
		Text:     quoteData.Text,
		Author:   quoteData.Author,
		Source:   "The Quote Hub",
		Language: "en",
	}
}

// GetLocationInformation performs a reverse geocoding request to get location information
func GetLocationInformation(latitude, longitude float64) (*db.LocationInfo, error) {
	url := fmt.Sprintf(
		"https://api.geoapify.com/v1/geocode/reverse?api_key=%s&lat=%f&lon=%f",
		os.Getenv("GEOAPIFY_API_KEY"),
		latitude,
		longitude,
	)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{
		Timeout: 30 * time.Second,
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	// Unmarshal response body to anonymous struct
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	// Define in-line struct for unmarshalling
	var locationData struct {
		Features []struct {
			Properties struct {
				Country string `json:"country"`
				City    string `json:"city"`
				State   string `json:"state"`
				Zip     string `json:"zip"`
				Lat     string `json:"lat"`
				Lon     string `json:"lon"`
			} `json:"properties"`
		} `json:"features"`
	}
	json.Unmarshal(resBody, &locationData)
	return &db.LocationInfo{
		Country: locationData.Features[0].Properties.Country,
		City:    locationData.Features[0].Properties.City,
		State:   locationData.Features[0].Properties.State,
		Zip:     locationData.Features[0].Properties.Zip,
		Lat:     locationData.Features[0].Properties.Lat,
		Lon:     locationData.Features[0].Properties.Lon,
	}, nil
}
