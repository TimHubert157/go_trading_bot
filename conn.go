package main

import (
	"log"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

type dataset struct {
	name string
	set  []float64
}

func WebsocketHandler(WaitForRequest chan bool) {
	conn, _, err := websocket.DefaultDialer.Dial("wss://stream.binance.com/ws", nil)
	requestBody, pairs := prepareRequest()
	symbols = pairs

	//prepare main dataset array
	for i := 0; i < len(pairs); i++ {
		n := dataset{}
		n.name = pairs[i]
		df = append(df, n)
	}

	defer wg.Done()

	if err != nil {
		log.Fatal(err)
	} else {
		defer conn.Close()

		werr := conn.WriteJSON(requestBody)
		if werr != nil {
			log.Fatal(werr)
		}
	}

	jsonBody := make(map[string]interface{})
	GotRequest := false

	for {
		rerr := conn.ReadJSON(&jsonBody)
		if rerr != nil {
			log.Fatal(rerr)
		}

		if jsonBody["k"] != nil {
			ClosingTime := time.UnixMilli(int64(jsonBody["k"].(map[string]interface{})["T"].(float64)))
			if time.Now().Add(15*time.Second).Before(ClosingTime) && GotRequest == false {
				getHistoricalKlines(pairs)
				WaitForRequest <- true
				GotRequest = true
			}
		}

		if jsonBody["k"] != nil && GotRequest {
			for i := 0; i < len(df); i++ {
				if df[i].name == jsonBody["s"] {
					// checking for new kline and appending a new one

					klines := jsonBody["k"].(map[string]interface{})

					if klines["x"].(bool) == true {
						floatAppend, _ := strconv.ParseFloat(klines["c"].(string), 32)
						df[i].set = append(df[i].set, floatAppend)
						// check for empty array, or to replace current value
					} else {
						floatReplace, _ := strconv.ParseFloat(klines["c"].(string), 32)
						if len(df[i].set) == 0 {
							df[i].set = append(df[i].set, floatReplace)
						} else {
							df[i].set[len(df[i].set)-1] = floatReplace
						}

					}
				}

			}

		}
	}

}
