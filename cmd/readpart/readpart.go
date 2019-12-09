package main

// Reads parititioned graph blobs

import (
	"fmt"
	"os"

	"github.com/sbromberger/gohavoq"
)

func main() {
	f := os.Args[1]
	for i := 0; i < 16; i++ {
		el2, _ := gohavoq.Load(fmt.Sprintf("%s-%d", f, i))  // el2 will be empty if file does not exist
		fmt.Println("-------------------------")
		fmt.Println("part ", i)
		fmt.Println(el2)
	}
}
