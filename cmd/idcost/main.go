package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		show_usage()
		return
	}
	data := os.Args[1]

	bytes := getBytes(data)
	zcount, nzcount, cost := getCost(bytes)

	fmt.Println("Non zero bytes: ", nzcount)
	fmt.Println("Zero bytes: ", zcount)
	fmt.Println("Total bytes: ", len(bytes))
	fmt.Println("Cost: ", cost)
}

func show_usage() {
	fmt.Println("idcost - calculates the cost for the given transaction input data")
	fmt.Println("Usage:")
	fmt.Println("\tidcost <txdata>")
}

func getBytes(data string) []string {
	if data[:2] == "0x" {
		data = data[2:]
	}
	totalBytes := len(data) / 2
	hexBytes := make([]string, totalBytes)

	i := 0
	j := 0
	for i < len(data) {
		hexBytes[j] = string(data[i]) + string(data[i+1])
		i += 2
		j += 1
	}

	return hexBytes
}

func getCost(bytes []string) (int, int, int) {
	cost := 0
	zero_count := 0
	non_zero_count := 0
	for _, b := range bytes {
		if b == "00" {
			cost += 4
			zero_count += 1
		} else {
			cost += 68
			non_zero_count += 1
		}
	}
	return zero_count, non_zero_count, cost
}
