package main

import (
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup
var df []dataset
var symbols []string

type bot struct {
	Symbol   string
	Bought   bool
	BoughtAt float64
	Dataset  *[]float64
}

func (b *bot) Strategy() {
	for {
		time.Sleep(1 * time.Second)
		fmt.Println(b.Symbol, (*b.Dataset)[len(*b.Dataset)-1])
	}
}

func main() {
	wg.Add(1)

	WaitForRequest := make(chan bool)

	go WebsocketHandler(WaitForRequest)

	if <-WaitForRequest {
		bots := make([]bot, 0)

		for i, symbol := range symbols {
			newBot := bot{symbol, false, 0, &df[i].set}
			wg.Add(1)
			go newBot.Strategy()
			bots = append(bots, newBot)
		}
	}

	wg.Wait()

}
