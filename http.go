package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

//Takes an array with symbols to get the Historical Klines (Only from Binance)
func getHistoricalKlines(Symbols []string) {

	for _, symbol := range Symbols {
		APIurl, err := url.Parse("https://api.binance.com/api/v3/klines")

		URLparams := url.Values{}

		timeNow := time.Now()
		URLparams.Add("symbol", symbol)
		URLparams.Add("interval", "5m")
		URLparams.Add("startTime", fmt.Sprintf("%d", timeNow.Add(-5*10*time.Minute).UnixMilli()))
		URLparams.Add("endTime", fmt.Sprintf("%d", timeNow.UnixMilli()))

		APIurl.RawQuery = URLparams.Encode()

		if err != nil {
			log.Fatal(err)
		}

		res, err := http.Get(APIurl.String())

		if err != nil {
			log.Fatal(err)
		}

		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}

		var ArrayRecv [][]interface{}

		jerr := json.Unmarshal(body, &ArrayRecv)
		if jerr != nil {
			log.Fatal(jerr)
		}

		for dfIndex, s := range symbols {
			if s == symbol {
				for _, oldKlines := range ArrayRecv {
					//convert json string to float for calc
					refloat, _ := strconv.ParseFloat(oldKlines[4].(string), 64)
					df[dfIndex].set = append(df[dfIndex].set, refloat)
				}
			}
		}

	}

}
