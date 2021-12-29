package h3geodist

import "github.com/uber/h3-go/v3"

// Iter iterate each cell at specified level, calling fn for each cell.
// Allowed levels 0-6.
func Iter(level int, fn func(index uint, cell h3.H3Index)) {
	if ok := validateLevel(level); !ok {
		return
	}
	var next uint
	if level == 0 {
		for _, cell0 := range h3.GetRes0Indexes() {
			next++
			fn(next, cell0)
		}
	} else {
		for _, cell0 := range h3.GetRes0Indexes() {
			for _, cell := range h3.ToChildren(cell0, level) {
				next++
				fn(next, cell)
			}
		}
	}
}

func validateLevel(level int) bool {
	if level < Level0 || level > Level6 {
		return false
	}
	return true
}
