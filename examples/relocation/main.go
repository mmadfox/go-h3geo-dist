package main

import (
	"fmt"

	h3geodist "github.com/mmadfox/go-h3geo-dist"
	"github.com/uber/h3-go/v3"
)

func main() {
	level := h3geodist.Level1
	area := h3geodist.Level1Area()

	h3dist, err := h3geodist.New(level)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 3; i++ {
		if err := h3dist.Add(fmt.Sprintf("127.0.0.%d", i)); err != nil {
			panic(err)
		}
	}

	// some owners
	owners := make(map[h3.H3Index]string)
	h3geodist.Iter(level, func(index uint, cell h3.H3Index) {
		dcell, ok := h3dist.Lookup(cell)
		if ok {
			owners[dcell.H3ID] = dcell.Host
		}
	})

	if err := h3dist.Add("127.0.0.4"); err != nil {
		panic(err)
	}

	stats := make(map[string]int)
	var changed int
	for cellID, oldhost := range owners {
		newhost, ok := h3dist.Lookup(cellID)
		if ok && newhost.Host != oldhost {
			changed++
			fmt.Printf("cellID: %v moved to %s from %s\n", cellID, newhost.Host, oldhost)
		}
		stats[newhost.Host]++
	}

	fmt.Printf("\n%d%% of the cells are relocated  changed=%d, total=%d\nstats:\n",
		(100*changed)/int(area), changed, area)

	for host, counter := range stats {
		fmt.Printf("host=%s, counter=%d\n", host, counter)
	}
}
