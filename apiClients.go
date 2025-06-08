package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// GetWeatherInfo returns a pointer to a WeatherInfo struct
// with weather information for the given city and country.
func GetWeatherInfo(city, country, units string) *WeatherInfo {
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
	// unmarshal response body to anonymous struct
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil
	}
	// define in-line struct for unmarshalling
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
	var weatherInfo WeatherInfo
	if units == "imperial" {
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
	} else if units == "metric" {
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
