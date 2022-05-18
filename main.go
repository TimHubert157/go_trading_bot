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
	Bulllock bool
	BoughtAt float64
	Dataset  *[]float64
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
}

func (b *bot) takeProfit() {
	if Percentage(b.BoughtAt, (*b.Dataset)[len(*b.Dataset)-1]) >= 0.2 {
		fmt.Printf("Sold with Profit (%s) %f %f \n", b.Symbol, b.BoughtAt, (*b.Dataset)[len(*b.Dataset)-1])
		b.Bought = false
		b.BoughtAt = 0
		b.Bulllock = true
	}
}

func (b *bot) stopLoss() {
	if Percentage(b.BoughtAt, (*b.Dataset)[len(*b.Dataset)-1]) <= -0.2 {
		fmt.Printf("Sold with Loss (%s) %f \n", b.Symbol, b.BoughtAt)
		b.Bought = false
		b.BoughtAt = 0
		b.Bulllock = true
	}
}

func main() {
	wg.Add(1)

	WaitForRequest := make(chan bool)

	go WebsocketHandler(WaitForRequest)

	bots := make([]bot, 0)

	if <-WaitForRequest {
		close(WaitForRequest)
		for i, symbol := range symbols {
			newBot := bot{symbol, false, true, 0, &df[i].set}
			wg.Add(1)
			go newBot.Strategy()
			bots = append(bots, newBot)
		}
	}

	wg.Wait()

}
