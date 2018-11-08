package main

import (
	"math"
	"fmt"
)

func main()  {
	var f1,f5 float64
	f5 = 11200
	f1 = 8100
	add := f5/f1
	fmt.Println(add)
	//tt := math.Pow(add, 0.25)
	tt := math.Round(math.Pow(f5/f1 , 0.25)*100)/100 - 1
	fmt.Println(tt)
}