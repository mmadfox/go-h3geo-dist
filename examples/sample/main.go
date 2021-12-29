package main

import (
	"fmt"

	h3geodist "github.com/mmadfox/go-h3geo-dist"
	"github.com/uber/h3-go/v3"
)

func main() {
	h3dist, err := h3geodist.New(h3geodist.Level5,
		h3geodist.WithVNodes(8),
	)
	if err != nil {
		panic(err)
	}

	_ = h3dist.Add("127.0.0.1")
	_ = h3dist.Add("127.0.0.2")

	cells, err := h3.HexRange(h3.FromString("854176affffffff"), 8)
	if err != nil {
		panic(err)
	}

	stats := make(map[string]int)
	for _, cell := range cells {
		level := h3.Resolution(cell)
		dcell, ok := h3dist.Lookup(cell)
		if ok {
			fmt.Printf("cell=%v, level=%d, host=%s\n", dcell.H3ID, level, dcell.Host)
			stats[dcell.Host]++
		}
	}

	fmt.Println("Stats:")
	for host, counter := range stats {
		fmt.Printf("- host=%s, counter=%d\n", host, counter)
	}
}
