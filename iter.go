package h3geodist

import "github.com/uber/h3-go/v3"

func Iter(level int, fn func(index uint, cell h3.H3Index)) {
	var index uint
	if level == 0 {
		for _, cell0 := range h3.GetRes0Indexes() {
			index++
			fn(index, cell0)
		}
	} else {
		for _, cell0 := range h3.GetRes0Indexes() {
			for _, cell := range h3.ToChildren(cell0, level) {
				index++
				fn(index, cell)
			}
		}
	}
}
