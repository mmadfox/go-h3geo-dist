package h3geodist

import (
	"fmt"
	"sort"
	"sync"

	"github.com/uber/h3-go/v3"
)

const (
	DefaultLevel  = 4
	DefaultVNodes = 32
)

type Distributed struct {
	mu       sync.RWMutex
	nodes    *rrw
	spans    []span
	level    int
	vNodes   uint
	numCells uint
}

type span struct {
	host  *node
	start uint64
	end   uint64
}

type node struct {
	addr   string
	weight int
}

func Default() *Distributed {
	dist, _ := New(DefaultLevel, DefaultVNodes)
	return dist
}

func New(cellLevel int, vNodes int) (*Distributed, error) {
	if cellLevel < Level0 && cellLevel > Level6 {
		return nil, fmt.Errorf("h3geodist: unknown level - got %d, expected [%d-%d]",
			cellLevel, Level0, Level6)
	}
	if vNodes < 3 {
		vNodes = 3
	}
	return &Distributed{
		nodes:    newrrw(),
		level:    cellLevel,
		vNodes:   uint(vNodes),
		numCells: cellArea(cellLevel),
	}, nil
}

func (d *Distributed) Add(addr string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.nodes.add(addr, 1)
	d.distribute()
}

func (d *Distributed) AddWithWeight(addr string, weight int) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.nodes.add(addr, weight)
	d.distribute()
}

func (d *Distributed) Remove(addr string) {

}

func (d *Distributed) Lookup(cell uint64) (addr string, ok bool) {
	if len(d.spans) == 0 {
		return
	}
	for i := 0; i < len(d.spans); i++ {
		if cell >= d.spans[i].start && cell < d.spans[i].end {
			addr = d.spans[i].host.addr
			ok = true
			break
		}
	}
	return
}

func (d *Distributed) distribute() {
	ring := d.fillspans()
	d.spans = make([]span, 0, d.vNodes)
	var end, start uint64
	for i := 0; i < len(ring); i++ {
		start = ring[i]
		if i+1 < len(ring) {
			end = ring[i+1]
		} else {
			end = start + 1
		}
		node := d.nodes.next()
		if node == nil {
			continue
		}
		d.spans = append(d.spans, span{
			start: start,
			end:   end,
			host:  node,
		})
	}
}

func (d *Distributed) fillspans() []uint64 {
	ring := make([]uint64, 0, d.vNodes)
	vnodes := d.vNodes * uint(d.nodes.size())
	step := d.numCells / vnodes
	nextIndex := step
	Iter(d.level, func(index uint, cell h3.H3Index) {
		if index == 1 {
			ring = append(ring, uint64(cell))
		}
		if index >= nextIndex {
			nextIndex = nextIndex + step
			ring = append(ring, uint64(cell))
		}
		if index == d.numCells {
			ring = append(ring, uint64(cell))
		}
	})
	sort.Slice(ring, func(i int, j int) bool {
		return ring[i] < ring[j]
	})
	return ring
}
