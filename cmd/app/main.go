// Package comment
package main

import (
	"fmt"

	methods "github.com/pavlosamtsov/lab1-tooling/internal"
)

func main() {
	methods.SayHelloTo("Pavlo")

	sum := methods.Sum([]int{15, 20, 25})
	fmt.Println(sum)

	sorted := methods.Sort([]int{5, 2, 8, 1, 9})
	fmt.Println(sorted)
}
