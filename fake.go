package main

import (
	"fmt"
)

func main() {
	arr1 := []int{1, 2, 3, 4, 5}
	myslice := arr1[1:4]
	des := make([]int, len(myslice))
	src := myslice
	copy(des, src)
	fmt.Println("Source array:", src)
	fmt.Println("Destination array:", des)
}
