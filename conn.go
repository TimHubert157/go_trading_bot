package main

import (
	"fmt"
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

	for i := 0; i < len(pairs); i++ {
		n := dataset{}
		n.name = pairs[i]
		df = append(df, n)
	}

	defer wg.Done()

	if err != nil {
		fmt.Println("WS ERROR")
	} else {
		defer conn.Close()

		werr := conn.WriteJSON(requestBody)
		if werr != nil {
			fmt.Println("WS WRITE ERROR")
		}
	}

	jsonBody := make(map[string]interface{})

	for {
		//second For Loop to avoid offsets.
		for {
			rerr := conn.ReadJSON(&jsonBody)
			if rerr != nil {
				fmt.Println("JSON READ ERROR")
			}
			if jsonBody["k"] != nil {
				ClosingTime := time.UnixMilli(int64(jsonBody["k"].(map[string]interface{})["T"].(float64)))
				if time.Now().Add(15 * time.Second).Before(ClosingTime) {
					getHistoricalKlines(pairs)
					WaitForRequest <- true
					break

				}
			}
		}
		rerr := conn.ReadJSON(&jsonBody)
		if rerr != nil {
			fmt.Println("JSON READ ERROR")
		} else {
			if jsonBody["k"] != nil {
				// checking for closing Time
				for i := 0; i < len(df); i++ {
					if df[i].name == jsonBody["s"] {
						// checking for new kline and appending a new one
						if jsonBody["k"].(map[string]interface{})["x"].(bool) == true {
							floatAppend, _ := strconv.ParseFloat(jsonBody["k"].(map[string]interface{})["c"].(string), 32)
							df[i].set = append(df[i].set, floatAppend)
							// check for empty array, or to replace current value
						} else {
							floatReplace, _ := strconv.ParseFloat(jsonBody["k"].(map[string]interface{})["c"].(string), 32)
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

}
