package scraper

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	pyScript string = "main.py"
	url      string = "https://www.worldometers.info/coronavirus/country/south-africa/"
)

type DataJSON struct {
	DailyCases []interface{} `json:"daily_cases"`
	Date       string        `json:"date"`
}

func getDailyCasesJS(doc *goquery.Document) ([]string, error) {
	var scriptTxt []string

	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		txt := s.Text()
		if strings.Contains(txt, "graph-cases-daily") {
			scriptTxt = strings.Split(txt, "\n")
		}
	})

	if len(scriptTxt) == 0 {
		return nil, errors.New("no script text found")
	}

	return scriptTxt, nil
}

func getJSArray(scriptTxt []string, prefix string) (string, error) {
	for _, line := range scriptTxt {
		if strings.Contains(line, prefix) {
			startIndex := strings.Index(line, "[")
			endIndex := strings.Index(line, "]")

			if startIndex == -1 || endIndex == -1 {
				return "", errors.New("invalid data format")
			}

			return line[startIndex+1 : endIndex], nil
		}
	}

	return "", errors.New("no data found")
}

func getCasesData(scriptTxt []string) ([]interface{}, error) {
	dataStr, err := getJSArray(scriptTxt, "data: [")
	if err != nil {
		return nil, err
	}

	dataSlice := strings.Split(dataStr, ",")

	data := make([]interface{}, len(dataSlice))

	for i, c := range dataSlice {
		if c == "null" {
			data[i] = nil
		} else {
			f, err := strconv.ParseFloat(c, 64)
			if err != nil {
				return nil, err
			}
			data[i] = f
		}
	}

	return data, nil
}

func getLastRecordDate(scriptTxt []string) (string, error) {
	dates, err := getJSArray(scriptTxt, "categories: [")

	if err != nil {
		return "", err
	}

	endIndex := len(dates) - 1
	startIndex := endIndex - 12

	date := dates[startIndex:endIndex]

	d, err := time.Parse("Jan 02, 2006", date)
	if err != nil {
		return "", err
	}

	return d.Format("2006-01-02"), nil
}

func generateCSV() error {
	cmd := exec.Command("python3", pyScript)

	_, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("Error in python script exec: %w", err)
	}

	return nil
}

func ScrapeAndParseData() error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return err
	}

	scriptTxt, err := getDailyCasesJS(doc)
	if err != nil {
		return err
	}

	cases, err := getCasesData(scriptTxt)
	if err != nil {
		return err
	}

	date, err := getLastRecordDate(scriptTxt)
	if err != nil {
		return err
	}

	data := DataJSON{
		DailyCases: cases,
		Date:       date,
	}

	file, err := os.Create("data.json")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(data)
	if err != nil {
		return err
	}

	return generateCSV()
}
