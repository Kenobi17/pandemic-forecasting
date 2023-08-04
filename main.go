package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Kenobi17/pandemic-forecasting/handler"
	"github.com/Kenobi17/pandemic-forecasting/scraper"
)

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

	http.HandleFunc("/api/forecast", handler.GetForecastHandler)
	http.HandleFunc("/api/forecast/fake", handler.GetFakeForecastHandler)
	http.HandleFunc("/api/forecast/download", handler.DownloadForecastHandler)
	http.HandleFunc("/api/forecast/fake/download", handler.DownloadFakeForecastHandler)
	http.HandleFunc("/", handler.ServeIndexFileHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("."))))

	log.Println("Starting server on port", getPort())
	log.Fatal(http.ListenAndServe(":"+getPort(), nil))
}
