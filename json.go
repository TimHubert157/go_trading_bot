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
}

// read config.json and prepare request for API
func prepareRequest() (Request request, pairs []string) {

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

		Request.ID = 1
		Request.Method = "SUBSCRIBE"
		Request.Params = params
	}

	return
}
