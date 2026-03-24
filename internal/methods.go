// Package methods package comment
package methods

import "fmt"

// SayHelloTo method comment
func SayHelloTo(name string) {
	fmt.Println("Hello, " + name + "!")
}

// Sum method comment
func Sum(numbers []int) int {
	sum := 0
	for _, value := range numbers {
		sum = sum + value
	}

	return sum
}

// Sort method comment
func Sort(numbers []int) []int {
	for i := 0; i < len(numbers)-1; i++ {
		for j := 0; j < len(numbers)-1-i; j++ {
			if numbers[j] > numbers[j+1] {
				numbers[j], numbers[j+1] = numbers[j+1], numbers[j]
			}
		}
	}
	return numbers
}
