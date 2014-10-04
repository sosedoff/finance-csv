package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

/*
Format description:
x: Stock exchange
s: Symbol
n: Name
j1: Market capitalization
s6: Revenue
l1: Last Trade (Price Only)
k: 52 Week High
j: 52 week Low
t8: 1 yr Target Price
*/

const (
	QuotesUrl    = "http://download.finance.yahoo.com/d/quotes.csv"
	QuotesFormat = "xsnj1s6l1kjt8"
	MaxSymbols   = 200
)

var allSymbols []string

func loadSymbols(path string) {
	body, err := ioutil.ReadFile(path)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	lines := strings.Split(strings.TrimSpace(string(body)), "\n")

	for _, line := range lines {
		items := strings.Split(line, "\t")
		allSymbols = append(allSymbols, items[0])
	}
}

func fetchCsvData(symbols []string) (string, error) {
	symbolsStr := strings.Join(symbols, ",")
	url := fmt.Sprintf("%s?f=%s&s=%s", QuotesUrl, QuotesFormat, symbolsStr)
	resp, err := http.Get(url)

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	csv := strings.TrimSpace(string(body))
	return csv, err
}

func main() {
	loadSymbols("./nasdaq.txt")
	loadSymbols("./nyse.txt")

	numPages := len(allSymbols) / MaxSymbols
	if numPages*MaxSymbols < len(allSymbols) {
		numPages++
	}

	for i := 0; i < numPages; i++ {
		chunk := allSymbols[MaxSymbols*i : MaxSymbols*(i+1)]
		csv, err := fetchCsvData(chunk)

		if err == nil {
			fmt.Println(csv)
			time.Sleep(100 * time.Millisecond)
		} else {
			time.Sleep(1 * time.Second)
		}
	}
}
