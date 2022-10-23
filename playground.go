package main

import "fmt"

type Point struct {
	x int
	y int
}

func main() {
	s := make(map[int][]Point)
	fmt.Println(s)
	a := Point{1, 1}
	s[1] = append(s[1], a)
	s[1] = append(s[1], a)
	fmt.Println(s)
}
