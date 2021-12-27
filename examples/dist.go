package main

import (
	"log"

	h3geodist "github.com/mmadfox/go-h3geo-dist"
	"github.com/uber/h3-go/v3"
)

func main() {
	level := h3geodist.Level1
	dist, err := h3geodist.New(level, h3geodist.DefaultVNodes)
	if err != nil {
		panic(err)
	}

	dist.Add("127.0.0.1")
	dist.Add("127.0.0.2")
	dist.Add("127.0.0.3")
	dist.AddWithWeight("127.0.0.4", 2)

	// iterate over all cells level one
	h3geodist.Iter(level, func(index uint, cell h3.H3Index) {
		host, ok := dist.Lookup(uint64(cell))
		log.Printf("cell=%d, host=%s, found=%v\n", uint64(cell), host, ok)
	})
}
