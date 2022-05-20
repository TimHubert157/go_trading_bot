package main

import (
	"fmt"
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
	requestBody, pairs, datasetSize, klineInterval := prepareRequest()

	fmt.Printf("\nStarting Crypto Bot with %v and Interval:%s \n\n", pairs, klineInterval)

	supIntervals := []string{"1m", "3m", "5m", "15m", "30m", "1h", "2h", "4h", "6h", "8h", "12h", "1d", "3d", "1w", "1M"}
	isSupported := false

	for _, elem := range supIntervals {
		if elem == klineInterval {
			isSupported = true
		}
	}

	if !isSupported {
		log.Fatal("Not supported interval. All intervals are supported: 1m, 3m, 5m, 15m, 30m, 1h, 2h, 4h, 6h, 8h, 12h, 1d, 3d, 1w, 1M")
	}

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
			if time.Now().Add(15*time.Second).Before(ClosingTime) && !GotRequest {
				getHistoricalKlines(pairs, datasetSize, klineInterval)
				WaitForRequest <- true
				GotRequest = true
			}
		}

		if jsonBody["k"] != nil && GotRequest {
			for i := 0; i < len(df); i++ {
				if df[i].name == jsonBody["s"] {
					if len(df[i].set) > datasetSize {
						df[i].set = df[i].set[1:len(df[i].set)]
					}
					// checking for new kline and appending a new one

					klines := jsonBody["k"].(map[string]interface{})

					if klines["x"].(bool) {
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
