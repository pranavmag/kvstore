package main

import "fmt"

func main() {
	var m map[string]int

	m["Key"] = 3

	v, ok := m["Key"]

	fmt.Println(v, ok)
}
