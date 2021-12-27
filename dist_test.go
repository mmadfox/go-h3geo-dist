package h3geodist

import (
	"fmt"
	"testing"

	"github.com/uber/h3-go/v3"
)

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
