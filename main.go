package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

var wg sync.WaitGroup
var df []dataset
var OpenTrades int64
var symbols []string
var stopLoss float64
var takeProfit float64
var BotArray []*bot
var Interval string

type bot struct {
	Symbol     string     `json:"symbol"`
	Bought     bool       `json:"bought"`
	Bulllock   bool       `json:"-"`
	BoughtAt   float64    `json:"-"`
	Profit     float64    `json:"profit"`
	Price      float64    `json:"price"`
	Dataset    *[]float64 `json:"-"`
	OpenTrades *int64     `json:"-"`
}

//Start strategy
func (b *bot) Strategy() {
	for {
		time.Sleep(1 * time.Second)

		b.Price = (*b.Dataset)[len(*b.Dataset)-1]

		//prevent the bot from buying too late.
		if EMA(*b.Dataset, 3) < EMA(*b.Dataset, 9) {
			b.Bulllock = false
		}

		if !b.Bought && !b.Bulllock {
			if EMA(*b.Dataset, 3) > EMA(*b.Dataset, 9) {
				b.Buy()
				b.OpenTrade()
			}
		}

		if b.Bought {
			b.Profit = Percentage(b.BoughtAt, (*b.Dataset)[len(*b.Dataset)-1])
			b.takeProfit()
			b.stopLoss()
		}
	}
}

func (b *bot) Buy() {
	b.BoughtAt = (*b.Dataset)[len(*b.Dataset)-1]
	b.Bought = true
	fmt.Printf("Bought %s at %f\n", b.Symbol, b.BoughtAt)
}

func (b *bot) takeProfit() {
	if Percentage(b.BoughtAt, (*b.Dataset)[len(*b.Dataset)-1]) >= takeProfit {
		fmt.Printf("Sold with Profit (%s) %f %f \n", b.Symbol, b.BoughtAt, (*b.Dataset)[len(*b.Dataset)-1])
		b.Bought = false
		b.BoughtAt = 0
		b.Bulllock = true
		b.Profit = 0
		b.CloseTrade()
	}
}

func (b *bot) stopLoss() {
	if Percentage(b.BoughtAt, (*b.Dataset)[len(*b.Dataset)-1]) <= (-1 * stopLoss) {
		fmt.Printf("Sold with Loss (%s) %f \n", b.Symbol, b.BoughtAt)
		b.Bought = false
		b.BoughtAt = 0
		b.Bulllock = true
		b.Profit = 0
		b.CloseTrade()

	}
}

func (b *bot) OpenTrade() {
	curr := *b.OpenTrades
	delta := 1
	curr++
	for {
		check := atomic.AddInt64(b.OpenTrades, int64(delta))
		if check == curr {
			break
		}
	}
}

func (b *bot) CloseTrade() {
	curr := *b.OpenTrades
	delta := -1
	curr--
	for {
		check := atomic.AddInt64(b.OpenTrades, int64(delta))
		if check == curr {
			break
		}
	}
}

func main() {
	wg.Add(1)

	WaitForRequest := make(chan bool)

	go WebsocketHandler(WaitForRequest)

	if <-WaitForRequest {
		close(WaitForRequest)
		wg.Add(1)

		for i, symbol := range symbols {
			newBot := bot{symbol, false, true, 0, 0, 0, &df[i].set, &OpenTrades}
			BotArray = append(BotArray, &newBot)
			wg.Add(1)
			go newBot.Strategy()
		}

	}

	wg.Wait()

}
