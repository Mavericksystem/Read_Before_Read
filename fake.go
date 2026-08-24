package main

import (
	"fmt"
)

func myfunc(x int, y int) int {
	return x + y
}

func main() {
	fmt.Println(myfunc(1, 2))
}
