package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type PointResponse struct {
	Properties struct {
		HourlyForecast string `json:"forecastHourly"`
	} `json:"properties"`
}

type CallWeatherGovReturn struct {
	HourlyProperty struct {
		Periods []struct {
			StartTime       string `json:"startTime"`
			EndTime         string `json:"endTime"`
			Temperature     int64  `json:"temperature"`
			TemperatureUnit string `json:"temperatureUnit"`
			ShortForecast   string `json:"shortForecast"`
		}
	} `json:"properties"`
}

func CallWeatherGov(lat, lng float64) (*CallWeatherGovReturn, error) {
	var pResponse PointResponse
	var weather CallWeatherGovReturn

	// Call the points endpoint at Weather.Gov to get the string we need for whatever location they are in
	resp, err := http.Get(fmt.Sprintf("https://api.weather.gov/points/%f,%f", lat, lng))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &pResponse)
	if err != nil {
		return nil, err
	}

	// Use the string from the points endpoint to get the hourly fourcast at that location
	res, err := http.Get(pResponse.Properties.HourlyForecast)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	bodyHourly, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bodyHourly, &weather)
	if err != nil {
		return nil, err
	}

	return &weather, nil
}
