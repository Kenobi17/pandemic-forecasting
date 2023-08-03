package main

import (
	"fmt"

	"github.com/Kenobi17/pandemic-forecasting/scraper"
)

func main() {

	d := scraper.GenForecastJSON()
	fmt.Println(d)
}
