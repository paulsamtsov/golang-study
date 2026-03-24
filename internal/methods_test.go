package methods

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"testing"
)

func TestSayHelloTo(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"regular name", "Pavlo", "Hello, Pavlo!\n"},
		{"empty name", "", "Hello, !\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, w, _ := os.Pipe()
			old := os.Stdout
			os.Stdout = w

			SayHelloTo(tt.input)

			//LINT error
			// w.Close()
			if err := w.Close(); err != nil {
				t.Errorf("failed to close pipe: %v", err)
			}
			os.Stdout = old

			out, _ := io.ReadAll(r)
			if string(out) != tt.expected {
				t.Errorf("got %q, want %q", string(out), tt.expected)
			}
		})
	}
}

func TestSum(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected int
	}{
		{"positive numbers", []int{1, 2, 3, 4}, 10},
		{"with zero", []int{0, 5, 10}, 15},
		{"single element", []int{42}, 42},
		{"empty slice", []int{}, 0},
		{"negative numbers", []int{-1, -2, -3}, -6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Sum(tt.input)
			if result != tt.expected {
				t.Errorf("got %d, want %d", result, tt.expected)
			}
		})
	}
}

func TestSort(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{"already sorted", []int{1, 2, 3}, []int{1, 2, 3}},
		{"reversed", []int{3, 2, 1}, []int{1, 2, 3}},
		{"unsorted", []int{5, 2, 8, 1}, []int{1, 2, 5, 8}},
		{"single element", []int{42}, []int{42}},
		{"empty slice", []int{}, []int{}},
		{"duplicates", []int{3, 1, 2, 1}, []int{1, 1, 2, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Sort(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("got %v, want %v", result, tt.expected)
			}
		})
	}
}

func ExampleSayHelloTo() {
	SayHelloTo("Pavlo")
	// Output: Hello, Pavlo!
}

func ExampleSum() {
	fmt.Println(Sum([]int{1, 2, 3}))
	// Output: 6
}

func ExampleSort() {
	fmt.Println(Sort([]int{3, 1, 2}))
	// Output: [1 2 3]
}
