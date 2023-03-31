package main

import "fmt"

func myFunc() int {
	a := 1
	b := 2
	
	if a < 2 {
		a = b
	}

	return a + b
}

func main() {
	fmt.Println(myFunc())
}