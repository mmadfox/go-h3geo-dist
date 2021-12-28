package main

import (
	"fmt"

	h3geodist "github.com/mmadfox/go-h3geo-dist"
	"github.com/uber/h3-go/v3"
)

type marker struct {
	cell h3geodist.Cell
	// some data
}

func main() {
	h3dist, err := h3geodist.New(h3geodist.Level1,
		h3geodist.WithVNodes(512),
	)
	if err != nil {
		panic(err)
	}

	target := h3.FromString("821fa7fffffffff")
	mymarker := marker{}

	if err := h3dist.Add("127.0.0.1"); err != nil {
		panic(err)
	}

	dcell, ok := h3dist.Lookup(target)
	if ok {
		mymarker.cell = dcell
	}

	if err := h3dist.Add("127.0.0.2"); err != nil {
		panic(err)
	}

	if !h3dist.IsOwned(mymarker.cell) {
		dcell, ok = h3dist.Lookup(target)
		fmt.Printf(" %v moved to %s from %s\n", dcell.H3ID, dcell.Host, mymarker.cell.Host)
	} else {
		fmt.Printf(" %v not moved \n", dcell.H3ID)
	}
}
