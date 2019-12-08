package main

import (
	"fmt"

	"github.com/sbromberger/gohavoq"
)

func main() {
	f := "elpart"
	for i := 0; i < 4; i++ {
		el2, _ := gohavoq.Load(fmt.Sprintf("%s-%d", f, i))
		fmt.Println("-------------------------")
		fmt.Println("part ", i)
		fmt.Println(el2)
	}
}
