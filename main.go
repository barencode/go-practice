package main

import "fmt"

func main() {
	balance := 0

	for i := 1; i <= 5; i++ {
		if i == 3 {
			continue
		}
		balance += 2
	}

	fmt.Println(balance)
}
