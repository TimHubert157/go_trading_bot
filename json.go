package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type request struct {
	Method string   `json:"method"`
	Params []string `json:"params"`
	ID     int      `json:"id"`
}

type config struct {
	Pairs         []string `json:"pairs"`
	KlineInterval string   `json:"klineInterval"`
	DatasetSize   int      `json:"datasetSize"`
	TakeProfit    float64  `json:"takeProfit"`
	StopLoss      float64  `json:"stopLoss"`
}

type wsJSON struct {
	Symbols    []string `json:"symbol"`
	Interval   string   `json:"interval"`
	OpenTrades int      `json:"openTrades"`
	Profit     float64  `json:"profit"`
	Bot        []bot    `json:"bots"`
}

var wsjson wsJSON

// read config.json and prepare request for API
func prepareRequest() (Request request, pairs []string, DatasetSize int, KlineInterval string) {

	configFile, oerr := os.Open("config.json")
	if oerr != nil {
		log.Fatal(oerr)
	} else {
		byteValue, _ := ioutil.ReadAll(configFile)
		jsonBody := config{}
		jerr := json.Unmarshal(byteValue, &jsonBody)

		if jerr != nil {
			log.Fatal(jerr)
			os.Exit(1)
		}

		var params []string

		for i := 0; i < len(jsonBody.Pairs); i++ {
			params = append(params, strings.ToLower(jsonBody.Pairs[i]+"@kline_"+jsonBody.KlineInterval))
			pairs = append(pairs, jsonBody.Pairs[i])
		}

		DatasetSize = jsonBody.DatasetSize
		KlineInterval = jsonBody.KlineInterval
		stopLoss = jsonBody.StopLoss
		takeProfit = jsonBody.TakeProfit

		wsjson.Interval = jsonBody.KlineInterval
		wsjson.Symbols = jsonBody.Pairs
		wsjson.OpenTrades = 0
		wsjson.Profit = 0.0
		wsjson.Bot = bots

		Request.ID = 1
		Request.Method = "SUBSCRIBE"
		Request.Params = params
	}

	return
}
