package h3geodist

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/uber/h3-go/v3"
)

func TestDistributed_AddRemove(t *testing.T) {
	t.Skip()
	level1cells := []uint64{
		uint64(h3.FromString("81823ffffffffff")),
		uint64(h3.FromString("8182bffffffffff")),
	}
	hosts := []string{
		"host-1.com",
		"host-2.com",
		"host-3.com",
	}
	h3dist, err := New(Level1, DefaultVNodes)
	if err != nil {
		t.Fatal(err)
	}

	rand.Seed(time.Now().UnixNano())

	var wg sync.WaitGroup
	wg.Add(2)
	const addOp = 0
	const removeOp = 1
	go func() {
		defer wg.Done()
		var index int
		for i := 0; i < 100; i++ {
			host := hosts[index%len(hosts)]
			switch rand.Intn(3) {
			case addOp:
				h3dist.Add(host)
			case removeOp:
				h3dist.Remove(host)
			}
			index++
			time.Sleep(15 * time.Millisecond)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 1300; i++ {
			cell := level1cells[i%2]
			host, ok := h3dist.Lookup(cell)
			t.Logf("cell=%d, host=%s, found=%v\n", uint64(cell), host, ok)
			time.Sleep(5 * time.Millisecond)
		}
	}()
	wg.Wait()
}

func TestDistributed_Remove(t *testing.T) {
	dist := Default()
	dist.Add("host-1.com")
	dist.Add("host-2.com")
	dist.Add("host-3.com")
	dist.Add("host-4.com")
	nodes := dist.Nodes()
	if have, want := len(nodes), 4; have != want {
		t.Fatalf("have %d, want %d", have, want)
	}

	dist.Remove("host-1.com")
	dist.Remove("host-2.com")
	dist.Remove("host-3.com")
	dist.Remove("host-4.com")
	nodes = dist.Nodes()
	if have, want := len(nodes), 0; have != want {
		t.Fatalf("have %d, want %d", have, want)
	}
}

func TestDistributed_Lookup(t *testing.T) {
	dist, _ := New(Level1, 9)
	dist.Add("host-1.com")
	dist.Add("host-2.com")
	dist.Add("host-3.com")
	dist.Add("host-4.com")
	level1Cells := []h3.H3Index{
		h3.FromString("8182bffffffffff"),
		h3.FromString("8158bffffffffff"),
		h3.FromString("81827ffffffffff"),
		h3.FromString("81547ffffffffff"),
		h3.FromString("817cfffffffffff"),
	}
	var found int
	for i := 0; i < len(level1Cells); i++ {
		host, ok := dist.Lookup(uint64(level1Cells[i]))
		t.Logf("host=%s, found=%v", host, ok)
		if ok {
			found++
		}
	}
	if have, want := len(level1Cells), found; have != want {
		t.Fatalf("have %d, want %d ", have, want)
	}
}

func TestDistributed_DynamicAddWithWeight(t *testing.T) {
	dist := Default()
	loop := 6
	for n := 0; n < loop; n++ {
		weight := 1
		if n > 2 && n < 4 {
			weight = n
		}
		addr := fmt.Sprintf("127.0.0.1:%d", n)
		dist.AddWithWeight(addr, weight)
		t.Logf("addNode=%s, weight=%d", addr, weight)

		var found uint
		stats := make(map[string]uint)
		Iter(DefaultLevel, func(i uint, cell h3.H3Index) {
			host, ok := dist.Lookup(uint64(cell))
			if ok {
				stats[host]++
				found++
			}
		})
		if have, want := found, Level4Area(); have != want {
			t.Fatalf("found %d, want %d cells", have, want)
		}
		var have uint
		for host, counter := range stats {
			have += counter
			t.Logf("host=%s, counter=%d", host, counter)
		}
		if found != have {
			t.Fatalf("have %d, want %d", have, found)
		}
	}
}

func TestDistributed_DynamicAdd(t *testing.T) {
	dist := Default()
	loop := 6
	for n := 0; n < loop; n++ {
		addr := fmt.Sprintf("127.0.0.1:%d", n)
		dist.Add(addr)
		t.Logf("addNode=%s", addr)

		var found uint
		stats := make(map[string]uint)
		Iter(DefaultLevel, func(i uint, cell h3.H3Index) {
			host, ok := dist.Lookup(uint64(cell))
			if ok {
				stats[host]++
				found++
			}
		})
		if have, want := found, Level4Area(); have != want {
			t.Fatalf("found %d, want %d cells", have, want)
		}
		var have uint
		for host, counter := range stats {
			have += counter
			t.Logf("host=%s, counter=%d", host, counter)
		}
		if found != have {
			t.Fatalf("have %d, want %d", have, found)
		}
	}
}
