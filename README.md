# H3-geo distributed cells 

[![Documentation](https://godoc.org/github.com/mmadfox/go-h3geo-dist?status.svg)](https://pkg.go.dev/github.com/mmadfox/go-h3geo-dist)
[![Coverage Status](https://coveralls.io/repos/github/mmadfox/go-h3geo-dist/badge.svg?branch=main)](https://coveralls.io/github/mmadfox/go-h3geo-dist?branch=main&1)
[![Go Report Card](https://goreportcard.com/badge/github.com/mmadfox/go-h3geo-dist)](https://goreportcard.com/report/github.com/mmadfox/go-h3geo-dist)

Distribution of [Uber H3geo](https://h3geo.org/) cells by nodes 

Prerequisites
-------
[H3-Go requires CGO ](https://github.com/uber/h3-go#prerequisites)

Install
-------
With a correctly configured Go env:

```
go get github.com/mmadfox/go-h3geo-dist
```

Examples
--------
```go
package main

import (
	"fmt"

	h3geodist "github.com/mmadfox/go-h3geo-dist"
	"github.com/uber/h3-go/v3"
)

func main() {
	level := h3geodist.Level1
	h3dist, err := h3geodist.New(level)
	if err != nil {
		panic(err)
	}

	_ = h3dist.Add("127.0.0.1")
	_ = h3dist.Add("127.0.0.2")
	_ = h3dist.Add("127.0.0.3")

	// iterate over all cells level one
	h3geodist.Iter(level, func(index uint, cell h3.H3Index) {
		// find a node by h3geo cell
		dcell, ok := h3dist.Lookup(cell)
		fmt.Printf("h3dist.Lookup: cell=%v, host=%s, found=%v\n", cell, dcell.Host, ok)
	})

	h3dist.LookupMany([]h3.H3Index{
		h3.FromString("821fa7fffffffff"),
		h3.FromString("821f9ffffffffff"),
		h3.FromString("81973ffffffffff"),
		h3.FromString("81f07ffffffffff"),
	}, func(c h3geodist.Cell) bool {
		fmt.Printf("h3dist.LookupMany: cell=%v, host=%s\n", c.H3ID, c.Host)
		return true
	})

	h3dist.Remove("127.0.0.1")
	h3dist.Remove("127.0.0.2")
	h3dist.Remove("127.0.0.3")
}
```