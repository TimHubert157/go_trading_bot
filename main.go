package main

import (
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup
var df []dataset
var bots []bot
var OpenTrades int
var symbols []string
var stopLoss float64
var takeProfit float64

type bot struct {
	Symbol     string
	Bought     bool
	Bulllock   bool
	BoughtAt   float64
	Dataset    *[]float64
	OpenTrades *int
}

//Start strategy
func (b *bot) Strategy() {
	for {
		time.Sleep(1 * time.Second)

		//prevent the bot from buying too late.
		if EMA(*b.Dataset, 3) < EMA(*b.Dataset, 9) {
			b.Bulllock = false
		}

		if !b.Bought && !b.Bulllock {
			if EMA(*b.Dataset, 3) > EMA(*b.Dataset, 9) {
				b.Buy()
			}
		}

		if b.Bought {
			b.takeProfit()
			b.stopLoss()
		}
	}
}

func (b *bot) Buy() {
	b.BoughtAt = (*b.Dataset)[len(*b.Dataset)-1]
	b.Bought = true
	fmt.Printf("Bought %s at %f\n", b.Symbol, b.BoughtAt)
	*b.OpenTrades++
}

func (b *bot) takeProfit() {
	if Percentage(b.BoughtAt, (*b.Dataset)[len(*b.Dataset)-1]) >= takeProfit {
		fmt.Printf("Sold with Profit (%s) %f %f \n", b.Symbol, b.BoughtAt, (*b.Dataset)[len(*b.Dataset)-1])
		b.Bought = false
		b.BoughtAt = 0
		b.Bulllock = true
		*b.OpenTrades--
	}
}

func (b *bot) stopLoss() {
	if Percentage(b.BoughtAt, (*b.Dataset)[len(*b.Dataset)-1]) <= (-1 * stopLoss) {
		fmt.Printf("Sold with Loss (%s) %f \n", b.Symbol, b.BoughtAt)
		b.Bought = false
		b.BoughtAt = 0
		b.Bulllock = true
		*b.OpenTrades--
	}
}

func main() {
	wg.Add(1)

	WaitForRequest := make(chan bool)

	go WebsocketHandler(WaitForRequest)

	if <-WaitForRequest {
		close(WaitForRequest)
		wg.Add(1)
		go Webserver()
		for i, symbol := range symbols {
			newBot := bot{symbol, false, true, 0, &df[i].set, &OpenTrades}
			wg.Add(1)
			go newBot.Strategy()
			bots = append(bots, newBot)
		}
	}

	wg.Wait()

}
