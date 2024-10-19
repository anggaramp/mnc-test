package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type StringCount struct {
	total int
	index []int
}

func main() {
	var n int
	fmt.Print("Input Number Of Array String: ")
	fmt.Scan(&n)

	arrayString := make([]string, n)

	reader := bufio.NewReader(os.Stdin)

	for i := 0; i < n; i++ {
		fmt.Printf("String %d: ", i+1)
		str, _ := reader.ReadString('\n')
		arrayString[i] = str[:len(str)-1]
	}

	stringCountMap := make(map[string]StringCount)

	//get total
	for i, v := range arrayString {
		if _, ok := stringCountMap[v]; ok {
			stringCountMap[v] = StringCount{
				total: stringCountMap[v].total + 1,
				index: append(stringCountMap[v].index, i+1),
			}
		} else {
			stringCountMap[v] = StringCount{
				total: 1,
				index: []int{i + 1},
			}
		}
	}
	var mostFrequentWord string
	maxTotal := 0

	for word, data := range stringCountMap {
		if data.total > maxTotal {
			maxTotal = data.total
			mostFrequentWord = word
		}
	}

	if maxTotal <= 1 {
		fmt.Println("Output : ", false)
		return
	}
	output := strings.Trim(strings.Join(strings.Split(fmt.Sprint(stringCountMap[mostFrequentWord].index), " "), " "), "[]")
	// Output the result
	fmt.Println("Output :", output)
}
