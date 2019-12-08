package main

import (
	"log"

	"github.com/sbromberger/gohavoq"
)

func partfn(u, v uint32, n int) (uint64, uint64, int) {
	return uint64(u), uint64(v), int(u) % n
}

func main() {
	f := "test.el"
	err := gohavoq.PartitionFromEdgeList(f, "elpart", partfn, 4, 0)
	if err != nil {
		log.Fatal("error: ", err)
	}
}
