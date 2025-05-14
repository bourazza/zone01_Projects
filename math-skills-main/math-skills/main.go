package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Print("zaki")
		return
	}

	var t int
	var sum float64
	var numbers []float64

	data := os.Args[1]
	file, err := os.Open(data)
	if err != nil {
		fmt.Print("err")
		return
	}
	defer file.Close()

	r := bufio.NewReader(file)

	for {
		line, _, err := r.ReadLine()

		if len(line) > 0 {
			f, err := strconv.ParseFloat(string(line), 64)
			if err != nil {
				fmt.Print(string(line), ": is not a float/int.  ")
				return
			}
			sum += f
			t++
			numbers = append(numbers, f)
		}

		if err != nil {
			break
		}
	}

	average := sum / float64(t)
	median := midian(numbers)
	variance := calcVariance(numbers, average)
	stdDev := math.Sqrt(variance)

	fmt.Println("average:", int(math.Round(average)))
	fmt.Println("median:", int(math.Round(median)))
	fmt.Println("variance:", int(math.Round(variance)))
	fmt.Println("standard deviation:", int(math.Round(stdDev)))
}

func midian(j []float64) float64 {
	sort.Float64s(j)
	n := len(j)
	if n == 0 {
		return 0
	}
	if n%2 == 1 {
		return j[n/2]
	}
	return (j[n/2-1] + j[n/2]) / 2
}

func calcVariance(nums []float64, mean float64) float64 {
	var sumSquares float64
	for _, v := range nums {
		diff := v - mean
		sumSquares += diff * diff
	}

	return sumSquares / float64(len(nums))
}
