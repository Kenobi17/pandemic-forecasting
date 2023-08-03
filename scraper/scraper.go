package scraper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const URL string = "https://www.worldometers.info/coronavirus/country/south-africa/"
const FilePath string = "data.json"

type DataJSON struct {
	DailyCases []interface{} `json:"daily_cases"`
	Date       string        `json:"date"`
}

func getScriptTxt(doc *goquery.Document) []string {
	var scriptTxt []string

	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		txt := s.Text()
		if strings.Contains(txt, "graph-cases-daily") {
			scriptTxt = strings.Split(txt, "\n")
		}
	})

	return scriptTxt
}

func getCasesData(scriptTxt []string) []interface{} {
	var dataStr string

	for _, line := range scriptTxt {
		if strings.Contains(line, "data: [") {
			dataStr = line

			startIndex := strings.Index(dataStr, "[")
			endIndex := strings.Index(dataStr, "]")

			dataStr = dataStr[startIndex+1 : endIndex]
			break
		}
	}

	dataSlice := strings.Split(dataStr, ",")

	data := make([]interface{}, len(dataSlice))

	for i, c := range dataSlice {
		if c == "null" {
			data[i] = nil
		} else {
			f, _ := strconv.ParseFloat(c, 64)
			data[i] = f
		}
	}

	return data
}

func getLastRecordDate(scriptTxt []string) string {
	var dateStr string

	for _, line := range scriptTxt {
		if strings.Contains(line, "categories: [") {
			dateStr = line

			endIndex := strings.Index(dateStr, "]")
			startIndex := endIndex - 13

			dateStr = dateStr[startIndex : endIndex-1]
			break
		}
	}

	date, _ := time.Parse("Jan 02, 2006", dateStr)

	return date.Format("2006-01-02")
}

func GenDailyCasesJSON() {
	res, err := http.Get(URL)

	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		panic(err)
	}

	scriptTxt := getScriptTxt(doc)

	cases := getCasesData(scriptTxt)

	date := getLastRecordDate(scriptTxt)

	d := DataJSON{
		DailyCases: cases,
		Date:       date,
	}

	file, _ := os.Create("data.json")
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.Encode(d)
}

func GenForecastJSON() string {
	_, err := os.Stat(FilePath)

	if os.IsNotExist(err) {
		GenDailyCasesJSON()
	}

	cmd := exec.Command("python3", "s.py")

	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error in python script exec", err)
	}

	return string(output)
}
