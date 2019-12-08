package main

import (
	"log"

	"github.com/sbromberger/gohavoq"
)

func main() {
	n := 3
	el := make(gohavoq.EdgeList, n)
	for i := uint64(0); i < uint64(n); i++ {
		e := gohavoq.Edge{Src: i, Dst: i + 1}
		el[i] = e
	}

	el.Save("test.el")

	el2, _ := gohavoq.Load("test.el")
	for i := range el {
		if el[i].Src != el2[i].Src || el[i].Dst != el2[i].Dst {
			log.Fatalf("%v != %v", el[i], el2[i])
		}
	}
}
