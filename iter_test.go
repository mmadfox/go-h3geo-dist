package h3geodist

import (
	"testing"

	"github.com/uber/h3-go/v3"
)

func TestIterLevel0(t *testing.T) {
	var cells uint
	Iter(Level0, func(index uint, cell h3.H3Index) {
		cells++
	})
	if have, want := cells, Level0Area(); have != want {
		t.Fatalf("have %d, want %d cells", have, want)
	}
}

func TestIterLevel1(t *testing.T) {
	var cells uint
	Iter(Level1, func(index uint, cell h3.H3Index) {
		cells++
	})
	if have, want := cells, Level1Area(); have != want {
		t.Fatalf("have %d, want %d cells", have, want)
	}
}

func TestIterLevel2(t *testing.T) {
	var cells uint
	Iter(Level2, func(index uint, cell h3.H3Index) {
		cells++
	})
	if have, want := cells, Level2Area(); have != want {
		t.Fatalf("have %d, want %d cells", have, want)
	}
}

func TestIterLevel3(t *testing.T) {
	var cells uint
	Iter(Level3, func(index uint, cell h3.H3Index) {
		cells++
	})
	if have, want := cells, Level3Area(); have != want {
		t.Fatalf("have %d, want %d cells", have, want)
	}
}

func TestIterLevel4(t *testing.T) {
	var cells uint
	Iter(Level4, func(index uint, cell h3.H3Index) {
		cells++
	})
	if have, want := cells, Level4Area(); have != want {
		t.Fatalf("have %d, want %d cells", have, want)
	}
}

func TestIterLevel5(t *testing.T) {
	var cells uint
	Iter(Level5, func(index uint, cell h3.H3Index) {
		cells++
	})
	if have, want := cells, Level5Area(); have != want {
		t.Fatalf("have %d, want %d cells", have, want)
	}
}

func TestIterLevel6(t *testing.T) {
	var cells uint
	Iter(Level6, func(index uint, cell h3.H3Index) {
		cells++
	})
	if have, want := cells, Level6Area(); have != want {
		t.Fatalf("have %d, want %d cells", have, want)
	}
}

func TestIterLevel7(t *testing.T) {
	var cells uint
	Iter(7, func(index uint, cell h3.H3Index) {
		cells++
	})
	if have, want := cells, uint(0); have != want {
		t.Fatalf("have %d, want %d cells", have, want)
	}
}
