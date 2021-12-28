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

func TestDistributed_LookupMany(t *testing.T) {
	h3dist, err := New(Level3,
		WithVNodes(9),
		WithReplicationFactor(9),
		WithLoadFactor(2),
	)
	if err != nil {
		t.Fatal(err)
	}

	for vn := uint64(0); vn < h3dist.VNodes(); vn++ {
		addr := fmt.Sprintf("host-%d.com", vn)
		if err := h3dist.Add(addr); err != nil {
			t.Fatal(err)
		}
	}

	var found int
	h3dist.LookupMany([]h3.H3Index{
		h3.FromString("821fa7fffffffff"),
		h3.FromString("821f9ffffffffff"),
		h3.FromString("81973ffffffffff"),
		h3.FromString("81f07ffffffffff"),
	}, func(c Cell) bool {
		found++
		return true
	})
	if have, want := found, 4; have != want {
		t.Fatalf("have %d, want %d num of cell", have, want)
	}
}

func TestDistributed_EachCell(t *testing.T) {
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
	h3dist.EachCell(func(c Cell) {
		stats[c.Host]++
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

func TestDistributed_IsOwned(t *testing.T) {
	h3dist, _ := New(Level1, WithVNodes(3))
	_ = h3dist.Add("127.0.0.1")
	want := h3.FromString("821fa7fffffffff")
	dcell, _ := h3dist.Lookup(want)
	if !h3dist.IsOwned(dcell) {
		t.Fatalf("h3dist.IsOwned(%v) => false, expected true", dcell)
	}
}

func TestDistributed_ReplicaFor(t *testing.T) {
	h3dist, _ := New(Level1, WithVNodes(256))
	_ = h3dist.Add("127.0.0.1")
	_ = h3dist.Add("127.0.0.2")
	_ = h3dist.Add("127.0.0.3")
	_ = h3dist.Add("127.0.0.4")
	cell := h3.FromString("821fa7fffffffff")
	hosts, err := h3dist.ReplicaFor(cell, 2)
	if err != nil {
		t.Fatal(err)
	}
	if have, want := len(hosts), 2; have != want {
		t.Fatalf("have %d, want %d", have, want)
	}
	_, err = h3dist.ReplicaFor(cell, 10)
	if err == nil {
		t.Fatalf("have nil, want error")
	}
}

func TestDefault(t *testing.T) {
	h3dist, _ := New(Level1, WithVNodes(3))

	want := h3.FromString("821fa7fffffffff")

	if err := h3dist.Add("127.0.0.1"); err != nil {
		t.Fatal(err)
	}

	dcell, ok := h3dist.Lookup(want)
	if ok {
		t.Logf("h3dist.Lookup(%v) => %s, %v", want, dcell.Host, ok)
	}

	if err := h3dist.Add("127.0.0.2"); err != nil {
		t.Fatal(err)
	}
	h3dist.Remove("127.0.0.1")

	dcell, ok = h3dist.Lookup(want)
	if ok {
		t.Logf("h3dist.Lookup(%v) => %s, %v", want, dcell.Host, ok)
	}
	h3dist.Remove("127.0.0.1")
	h3dist.Remove("127.0.0.2")
	dcell, ok = h3dist.Lookup(want)
	if !ok {
		t.Logf("h3dist.Lookup(%v) => %v", want, ok)
	}
}
