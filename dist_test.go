package h3geodist

import (
	"errors"
	"fmt"
	"testing"

	"github.com/uber/h3-go/v3"
)

func TestDistributed_Add(t *testing.T) {
	h3dist, err := New(Level3)
	if err != nil {
		t.Fatal(err)
	}

	for vn := uint64(0); vn < h3dist.VNodes(); vn++ {
		addr := fmt.Sprintf("host-%d.com", vn)
		if err := h3dist.Add(addr); err != nil {
			t.Fatal(err)
		}
	}

	stats := make(map[string]int)
	Iter(Level3, func(index uint, cell h3.H3Index) {
		dcell, ok := h3dist.Lookup(cell)
		if ok {
			stats[dcell.Host]++
		}
	})
	var total int
	for host, counter := range stats {
		t.Logf("host=%s, counter=%d", host, counter)
		total += counter
	}

	if have, want := uint(total), Level3Area(); have != want {
		t.Fatalf("have %d, want %d num of cell", have, want)
	}
}

func TestDistributed_AddWithError(t *testing.T) {
	h3dist, err := New(Level3)
	if err != nil {
		t.Fatal(err)
	}

	for i := uint64(0); i < h3dist.VNodes(); i++ {
		_ = h3dist.Add(fmt.Sprintf("127.0.0.%d", i))
	}

	if err := h3dist.Add("127.0.1.1"); !errors.Is(err, ErrNoSlots) {
		t.Fatalf("have %v, want %v error", err, ErrNoSlots)
	}
}

func TestDistributed_Remove(t *testing.T) {
	h3dist, err := New(Level3)
	if err != nil {
		t.Fatal(err)
	}
	for i := uint64(0); i < h3dist.VNodes()/2; i++ {
		if err := h3dist.Add(fmt.Sprintf("127.0.0.%d", i)); err != nil {
			t.Fatal(err)
		}
	}
	nodes := h3dist.Nodes()
	if have, want := len(nodes), int(h3dist.VNodes()/2); have != want {
		t.Fatalf("have %d, want %d", have, want)
	}

	for i := uint64(0); i < h3dist.VNodes()/2; i++ {
		h3dist.Remove(fmt.Sprintf("127.0.0.%d", i))
	}

	nodes = h3dist.Nodes()
	if have, want := len(nodes), 0; have != want {
		t.Fatalf("have %d, want %d", have, want)
	}
}
