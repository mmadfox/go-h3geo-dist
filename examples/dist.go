package main

import (
	"log"

	h3geodist "github.com/mmadfox/go-h3geo-dist"
	"github.com/uber/h3-go/v3"
)

func main() {
	level := h3geodist.Level1
	h3dist, err := h3geodist.New(level, h3geodist.DefaultVNodes)
	if err != nil {
		panic(err)
	}

	h3dist.Add("127.0.0.1")
	h3dist.Add("127.0.0.2")
	h3dist.Add("127.0.0.3")
	h3dist.AddWithWeight("127.0.0.4", 2)

	nodes := h3dist.Nodes()
	for i := 0; i < len(nodes); i++ {
		log.Printf("host=%s\n", nodes[i])
	}

	// iterate over all cells level one
	h3geodist.Iter(level, func(index uint, cell h3.H3Index) {
		// find a node by h3geo cell
		host, ok := h3dist.Lookup(uint64(cell))
		log.Printf("cell=%d, host=%s, found=%v\n", uint64(cell), host, ok)
	})

	h3dist.Remove("127.0.0.1")
	h3dist.Remove("127.0.0.2")
	h3dist.Remove("127.0.0.3")
	h3dist.Remove("127.0.0.4")
}
