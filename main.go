package main

import (
	"encoding/csv"
	"encoding/json"
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
const fakeForecastPath string = "fake_forecast.csv"

func getForecastFromCSV(w http.ResponseWriter, path string) {
	file, err := os.Open(path)
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
}

func getPort() string {
	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = "8080"
	}
	return port
}

func main() {
	err := scraper.ScrapeAndParseData()
	if err != nil {
		log.Fatalf("Error generating forecast data: %v", err)
	}

	http.HandleFunc("/api/forecast", func(w http.ResponseWriter, r *http.Request) {
		getForecastFromCSV(w, forecastPath)
	})

	http.HandleFunc("/api/forecast/fake", func(w http.ResponseWriter, r *http.Request) {
		getForecastFromCSV(w, fakeForecastPath)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("."))))

	log.Println("Starting server on port", getPort())
	log.Fatal(http.ListenAndServe(":"+getPort(), nil))
}
