package service

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Env struct {
	Port string
	Host string
}

type Handler struct {
	*Env
	H func(e *Env, w http.ResponseWriter, r *http.Request) error
}

type kv struct {
	Key   string
	Value int
}

// ServeHTTP allows our Handler type to satisfy http.Handler.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.H(h.Env, w, r)
	if err != nil {
		switch e := err.(type) {
		case Error:
			log.Printf("HTTP %d - %s", e.Status(), e)
			http.Error(w, e.Error(), e.Status())
		default:
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
	}
}

// GetWeather will receive lat and long within the http URL params and return weather of that location
func GetWeather(env *Env, w http.ResponseWriter, r *http.Request) error {
	lat, err := strconv.ParseFloat(strings.TrimSpace(r.URL.Query().Get("lat")), 64)
	if err != nil {
		return StatusError{500, err}
	}

	lng, err := strconv.ParseFloat(strings.TrimSpace(r.URL.Query().Get("lng")), 64)
	if err != nil {
		return StatusError{500, err}
	}

	res, err := CallWeatherGov(lat, lng)
	if err != nil {
		return StatusError{500, err}
	}

	var avgTemp float64
	shortForecast := make(map[string]int)

	// Go through all the hourly forecasts and only pick out todays'
	for i, hour := range res.HourlyProperty.Periods {
		parsedTime, err := time.Parse(time.RFC3339, hour.StartTime)
		if err != nil {
			return StatusError{500, err}
		}

		// Check to make sure we are still in the same days forecast
		if parsedTime.Day() != time.Now().Day() {
			avgTemp = avgTemp / float64(i+1)
			break
		}

		_, ok := shortForecast[hour.ShortForecast]
		if ok {
			shortForecast[hour.ShortForecast] = shortForecast[hour.ShortForecast] + 1
		} else {
			shortForecast[hour.ShortForecast] = 1
		}

		avgTemp += float64(hour.Temperature)
	}

	temperature := "cold"
	if math.Abs(avgTemp) > 50 {
		temperature = "moderate"
	}
	if math.Abs(avgTemp) > 90 {
		temperature = "hot"
	}

	w.Write([]byte(fmt.Sprintln("\nThe day will be", strings.ToLower(mostCommonForecast(shortForecast)), "with a", temperature, "average temperature.")))
	return nil
}

// Return the most common forecast of all the ones received for that day
func mostCommonForecast(m map[string]int) string {
	var ss []kv
	for k, v := range m {
		ss = append(ss, kv{k, v})
	}

	// Then sorting the slice by value, higher first.
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	return ss[0].Key
}
