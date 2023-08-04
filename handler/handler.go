package handler

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
)

type Forecast struct {
	DailyCases []int    `json:"daily_cases"`
	Dates      []string `json:"dates"`
}

const (
	forecastPath     = "forecast.csv"
	fakeForecastPath = "fake_forecast.csv"
)

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

func downloadCSVFile(w http.ResponseWriter, path string) {
	file, err := os.Open(path)
	if err != nil {
		http.Error(w, "Could not fetch data", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename="+path)
	io.Copy(w, file)
}

func GetForecastHandler(w http.ResponseWriter, r *http.Request) {
	getForecastFromCSV(w, forecastPath)
}

func GetFakeForecastHandler(w http.ResponseWriter, r *http.Request) {
	getForecastFromCSV(w, fakeForecastPath)
}

func DownloadForecastHandler(w http.ResponseWriter, r *http.Request) {
	downloadCSVFile(w, forecastPath)
}

func DownloadFakeForecastHandler(w http.ResponseWriter, r *http.Request) {
	downloadCSVFile(w, fakeForecastPath)
}

func ServeIndexFileHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}
