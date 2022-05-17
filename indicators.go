package main

// Simple Moving Average
func sma(arr []float64, period int) (sma float64) {
	ArraySma := []float64{}
	for i := len(arr); i > period; i-- {
		toSum := arr[i-period : i]
		var sum float64
		for _, elem := range toSum {
			sum = sum + float64(elem)
		}
		ArraySma = append(ArraySma, sum/float64(period))
		sum = 0
	}

	sma = ArraySma[0]

	return
}
