package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
		fmt.Println("FILE ERROR")
	} else {
		byteValue, _ := ioutil.ReadAll(configFile)
		jsonBody := config{}
		json.Unmarshal(byteValue, &jsonBody)

		var params []string

		for i := 0; i < len(jsonBody.Pairs); i++ {
			params = append(params, strings.ToLower(jsonBody.Pairs[i]+"@kline_"+jsonBody.KlineInterval))
			pairs = append(pairs, jsonBody.Pairs[i])
		}

		Request.ID, Request.Method, Request.Params = 1, "SUBSCRIBE", params
	}

	return
}
