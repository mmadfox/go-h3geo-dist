package h3geodist

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"strconv"
	"sync"

	"github.com/uber/h3-go/v3"
)

var ErrNoSlots = errors.New("h3geodist: no distribute slots")

type Distributed struct {
	mu         sync.RWMutex
	replFactor int
	loadFactor float64
	vnodes     uint64
	ring       map[uint64]*node
	index      map[int]*node
	hashes     []uint64
	nodes      []*node
	level      int
	stats      map[string]float64
}

type Cell struct {
	H3ID h3.H3Index
	Host string
}

type node struct {
	addr string
}

func Default() *Distributed {
	dist, _ := New(DefaultLevel)
	return dist
}

func New(cellLevel int, opts ...Option) (*Distributed, error) {
	if ok := validateLevel(cellLevel); !ok {
		return nil, fmt.Errorf("h3geodist: unsupported level - got %d, expected [%d-%d]",
			cellLevel, Level0, Level6)
	}
	h3dist := &Distributed{
		loadFactor: DefaultLoadFactor,
		replFactor: DefaultReplicationFactor,
		vnodes:     DefaultVNodes,
		level:      cellLevel,
		ring:       make(map[uint64]*node),
		index:      make(map[int]*node),
		stats:      make(map[string]float64),
	}
	for _, f := range opts {
		f(h3dist)
	}
	return h3dist, nil
}

func (d *Distributed) IsEmpty() bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return len(d.nodes) == 0
}

func (d *Distributed) VNodes() uint64 {
	return d.vnodes
}

func (d *Distributed) Nodes() []string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	nodes := make([]string, 0, len(d.nodes))
	for i := 0; i < len(d.nodes); i++ {
		nodes = append(nodes, d.nodes[i].addr)
	}
	return nodes
}

func (d *Distributed) Lookup(cell h3.H3Index) (Cell, bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if len(d.nodes) == 0 {
		return Cell{}, false
	}
	addr, ok := d.lookup(cell)
	if !ok {
		return Cell{}, false
	}
	return Cell{H3ID: cell, Host: addr}, true
}

func (d *Distributed) IsOwned(c Cell) bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	addr, ok := d.lookup(c.H3ID)
	if !ok {
		return false
	}
	return addr == c.Host
}

func (d *Distributed) LookupMany(cell []h3.H3Index, iter func(c Cell) bool) bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if len(cell) == 0 || len(d.nodes) == 0 {
		return false
	}
	for i := 0; i < len(cell); i++ {
		addr, ok := d.lookup(cell[i])
		if !ok {
			continue
		}
		if ok := iter(Cell{H3ID: cell[i], Host: addr}); !ok {
			return false
		}
	}
	return true
}

func (d *Distributed) EachCell(iter func(c Cell)) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if len(d.nodes) == 0 {
		return
	}
	Iter(d.level, func(_ uint, cell h3.H3Index) {
		addr, ok := d.lookup(cell)
		if !ok {
			return
		}
		iter(Cell{H3ID: cell, Host: addr})
	})
}

func (d *Distributed) Add(addr string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.exist(addr) {
		return nil
	}
	newNode := &node{addr: addr}
	d.nodes = append(d.nodes, newNode)
	d.add(newNode)
	return d.distribute()
}

func (d *Distributed) Remove(addr string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if !d.exist(addr) {
		return
	}
	d.remove(addr)
	_ = d.distribute()
}

func (d *Distributed) lookup(cell h3.H3Index) (addr string, ok bool) {
	hashKey := uint2hash(uint64(cell))
	idx := int(hashKey % d.vnodes)
	node, found := d.index[idx]
	if !found {
		return
	}
	addr = node.addr
	ok = true
	return
}

func (d *Distributed) exist(addr string) (ok bool) {
	for i := 0; i < len(d.nodes); i++ {
		if addr == d.nodes[i].addr {
			ok = true
			break
		}
	}
	return
}

func (d *Distributed) distribute() error {
	stats := make(map[string]float64)
	index := make(map[int]*node)
	for vnode := uint64(0); vnode < d.vnodes; vnode++ {
		nodeIndex := d.findNodeIndex(uint2hash(vnode))
		avgload := d.AvgLoad()
		var next int
		for {
			next++
			if next >= len(d.hashes) {
				return ErrNoSlots
			}
			node := d.ring[d.hashes[nodeIndex]]
			load := stats[node.addr]
			if load+1 <= avgload {
				index[int(vnode)] = node
				stats[node.addr]++
				break
			}
			nodeIndex++
			if nodeIndex >= len(d.hashes) {
				nodeIndex = 0
			}
		}
	}
	d.index = index
	d.stats = stats
	return nil
}

func (d *Distributed) AvgLoad() float64 {
	if len(d.nodes) == 0 {
		return 0
	}
	return math.Ceil(float64(d.vnodes/uint64(len(d.nodes))) * d.loadFactor)
}

func (d *Distributed) findNodeIndex(hashKey uint64) int {
	nodeIndex := sort.Search(len(d.hashes), func(n int) bool {
		return d.hashes[n] >= hashKey
	})
	if nodeIndex >= len(d.hashes) {
		nodeIndex = 0
	}
	return nodeIndex
}

func (d *Distributed) add(n *node) {
	for i := 0; i < d.replFactor; i++ {
		hashKey := str2hash(n.addr + strconv.Itoa(i))
		d.ring[hashKey] = n
		d.hashes = append(d.hashes, hashKey)
	}
	sort.Slice(d.hashes, func(i int, j int) bool {
		return d.hashes[i] < d.hashes[j]
	})
}

func (d *Distributed) remove(addr string) {
	for i := 0; i < d.replFactor; i++ {
		hashKey := str2hash(addr + strconv.Itoa(i))
		delete(d.ring, hashKey)
		for i := 0; i < len(d.hashes); i++ {
			if d.hashes[i] == hashKey {
				d.hashes = append(d.hashes[:i], d.hashes[i+1:]...)
				break
			}
		}
	}
	for i := 0; i < len(d.nodes); i++ {
		if d.nodes[i].addr == addr {
			d.nodes = append(d.nodes[:i], d.nodes[i+1:]...)
		}
	}
	delete(d.stats, addr)
	if len(d.nodes) == 0 {
		d.stats = make(map[string]float64)
		d.index = make(map[int]*node)
		d.hashes = make([]uint64, 0, d.vnodes)
	}
}
