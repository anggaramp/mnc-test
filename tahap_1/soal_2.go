package main

import (
	"fmt"
)

func main() {
	var totalPayment, userPayment int
	fmt.Print("Total belanja seorang customer : ")
	fmt.Scan(&totalPayment)
	fmt.Print("Pembeli membayar               : ")
	fmt.Scan(&userPayment)

	if userPayment < totalPayment {
		fmt.Println(false, "kurang bayar")
		return
	}

	returnPayment := userPayment - totalPayment

	fmt.Printf("Kembalian yang harus diberikan kasir: %d\n", returnPayment)

	returnPayment = returnPayment - (returnPayment % 100)

	fmt.Printf("Dibulatkan menjadi: %d\n", returnPayment)

	// Breakdown of denominations
	currency := []int{50000, 20000, 5000, 2000, 200, 100}
	quantity := make([]int, len(currency))

	for i, curr := range currency {
		if returnPayment >= curr {
			quantity[i] = returnPayment / curr
			returnPayment %= curr
		}
	}

	// Output the rounded change and denominations
	fmt.Println("Pecahan uang:")
	for i, curr := range currency {
		if quantity[i] > 0 {
			if curr >= 1000 {
				fmt.Printf("%d lembar %d\n", quantity[i], curr)
			} else {
				fmt.Printf("%d koin %d\n", quantity[i], curr)
			}
		}
	}
}
