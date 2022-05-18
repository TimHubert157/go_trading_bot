package main

// Simple Moving Average
func SMA(arr []float64, period int) (sma float64) {
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

//EMA today = α Price today + (1 − α) EMA yesterday | α = 2 / (n + 1) n = period
func EMA(arr []float64, period int) (ema float64) {
	ArrayEma := []float64{}
	a := 2 / (float64(period) + 1)
	for index, elem := range arr {
		if index == 0 {
			ArrayEma = append(ArrayEma, elem)
		} else {
			EMAtoday := (elem * a) + (1-a)*(ArrayEma[index-1])
			ArrayEma = append(ArrayEma, EMAtoday)
		}
	}

	return ArrayEma[len(ArrayEma)-1]
}

//Percentage difference between to values
func Percentage(s float64, f float64) (diff float64) {
	diff = (((s - f) / ((s + f) / 2) * 100) * -1)
	return
}
