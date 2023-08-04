package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/Kenobi17/pandemic-forecasting/scraper"
)

type Forecast struct {
	DailyCases []int    `json:"daily_cases"`
	Dates      []string `json:"dates"`
}

const forecastPath string = "forecast.csv"
const mockForecastPath string = "mock_forecast.csv"

func main() {
	err := scraper.ScrapeAndParseData()
	if err != nil {
		log.Fatalf("Error generating forecast data: %v", err)
	}

	http.HandleFunc("/api/forecast", func(w http.ResponseWriter, r *http.Request) {
		file, err := os.Open(forecastPath)
		if err != nil {
			http.Error(w, "Could not fetch forecast data", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		reader := csv.NewReader(file)
		rows, err := reader.ReadAll()
		if err != nil {
			http.Error(w, "Could not parse forecast data", http.StatusInternalServerError)
			return
		}

		forecast := Forecast{
			DailyCases: make([]int, len(rows)-1),
			Dates:      make([]string, len(rows)-1),
		}

		for i, r := range rows[1:] {
			forecast.Dates[i] = r[0]
			forecast.DailyCases[i], _ = strconv.Atoi(r[1])
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(forecast)
	})

	http.HandleFunc("/api/forecast/mock", func(w http.ResponseWriter, r *http.Request) {
		file, err := os.Open(mockForecastPath)
		if err != nil {
			http.Error(w, "Could not fetch forecast data", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		reader := csv.NewReader(file)
		rows, err := reader.ReadAll()
		if err != nil {
			http.Error(w, "Could not parse forecast data", http.StatusInternalServerError)
			return
		}

		forecast := Forecast{
			DailyCases: make([]int, len(rows)-1),
			Dates:      make([]string, len(rows)-1),
		}

		for i, r := range rows[1:] {
			forecast.Dates[i] = r[0]
			forecast.DailyCases[i], _ = strconv.Atoi(r[1])
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(forecast)
	})

	fmt.Println("Starting server on port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
