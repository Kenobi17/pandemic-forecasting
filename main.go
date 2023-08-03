package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getDailyCases() string {
	res, err := http.Get("https://www.worldometers.info/coronavirus/country/south-africa/")
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		panic(err)
	}

	var data string
	doc.Find("script").Each(func(i int, s *goquery.Selection) {

		txt := s.Text()
		lines := strings.Split(txt, "\n")

		for _, l := range lines {
			if strings.Contains(l, "data: [") {
				data = l

				startIndex := strings.Index(data, "[")
				endIndex := strings.Index(data, "]")

				data = data[startIndex : endIndex+1]
				break
			}
		}

	})

	return data
}

func main() {

	d := getDailyCases()

	os.WriteFile("data.json", []byte(d), 0644)

	cmd := exec.Command("python3", "s.py")

	output, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	fmt.Println(string(output))
}
