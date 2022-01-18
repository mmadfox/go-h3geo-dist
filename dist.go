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

var (
	// ErrNoSlots means that the number of virtual nodes is distributed by 100%.
	// It is necessary to change the configuration of virtual nodes.
	ErrNoSlots = errors.New("h3geodist: no distribute slots")

	// ErrVNodes returns when there are no virtual nodes.
	ErrVNodes = errors.New("h3geodist: vnodes not found")
)

// Distributed holds information about nodes,
// and scheduler of virtual nodes with replicas.
// Thread-safe.
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

// Cell is a type to represent a distributed cell
// with specifying the hostname and H3 Index.
type Cell struct {
	H3ID h3.H3Index
	Host string
}

func (c Cell) String() string {
	return fmt.Sprintf("Cell{Host: %s, ID: %s}", c.Host, h3.ToString(c.H3ID))
}

func (c Cell) HexID() string {
	return h3.ToString(c.H3ID)
}

// NodeInfo is a type to represent a node load statistic.
type NodeInfo struct {
	Host string
	Load float64
}

type node struct {
	addr string
}

// Default creates and returns a new Distributed instance with level - Level5.
func Default() *Distributed {
	dist, _ := New(Level5)
	return dist
}

// New creates and returns a new Distributed instance
// with specified cell level and options.
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

// IsEmpty returns TRUE if the nodes list are empty, otherwise FALSE.
func (d *Distributed) IsEmpty() bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return len(d.nodes) == 0
}

// NumReplica returns number of replicas.
func (d *Distributed) NumReplica() int {
	return d.replFactor
}

// Stats returns load distribution by nodes.
func (d *Distributed) Stats() []NodeInfo {
	d.mu.RLock()
	defer d.mu.RUnlock()
	stats := make([]NodeInfo, 0, len(d.stats))
	for host, load := range d.stats {
		stats = append(stats, NodeInfo{
			Host: host,
			Load: load,
		})
	}
	return stats
}

// VNodes returns number of virtual nodes.
func (d *Distributed) VNodes() uint64 {
	return d.vnodes
}

// Nodes returns a list of nodes.
func (d *Distributed) Nodes() []string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	nodes := make([]string, 0, len(d.nodes))
	for i := 0; i < len(d.nodes); i++ {
		nodes = append(nodes, d.nodes[i].addr)
	}
	return nodes
}

// Lookup returns distributed cell.
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

// IsOwned сhecks if the host for a distributed cell has changed.
func (d *Distributed) IsOwned(c Cell) bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	addr, ok := d.lookup(c.H3ID)
	if !ok {
		return false
	}
	return addr == c.Host
}

// WhereIsMyParent finds and returns parent distributed cell.
// The child object must be less resolution than the parent's parent.
func (d *Distributed) WhereIsMyParent(child h3.H3Index) (c Cell, err error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	curLevel := h3.Resolution(child)
	if curLevel < d.level {
		return c, fmt.Errorf("h3geodist: child resolution got %d, expected > %d",
			curLevel, d.level)
	}
	cell := h3.ToParent(child, d.level)
	addr, ok := d.lookup(cell)
	if !ok {
		return c, ErrVNodes
	}
	c.H3ID = cell
	c.Host = addr
	return
}

// LookupFromLatLon returns distributed cell.
func (d *Distributed) LookupFromLatLon(lat float64, lon float64) (c Cell, err error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	cell := h3.FromGeo(h3.GeoCoord{Latitude: lat, Longitude: lon}, d.level)
	addr, ok := d.lookup(cell)
	if !ok {
		return c, ErrVNodes
	}
	return Cell{H3ID: cell, Host: addr}, nil
}

// Neighbor is a type for represent a neighbor distributed cell,
// with the distance from a target point to the center of each neighbor.
type Neighbor struct {
	Cell      Cell
	DistanceM float64
}

// NeighborsFromLatLon returns the current distributed cell
// for a geographic coordinate and neighbors sorted by distance in descending order.
// Distance is measured from geographic coordinates to the center of each neighbor.
func (d *Distributed) NeighborsFromLatLon(lat float64, lon float64) (target Cell, neighbors []Neighbor, err error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	src := h3.GeoCoord{Latitude: lat, Longitude: lon}
	cell := h3.FromGeo(src, d.level)
	addr, ok := d.lookup(cell)
	if !ok {
		return target, nil, ErrVNodes
	}
	target.Host = addr
	target.H3ID = cell
	ring := h3.KRing(cell, 1)
	neighbors = make([]Neighbor, 0, len(ring))
	for i := 0; i < len(ring); i++ {
		if !h3.AreNeighbors(cell, ring[i]) {
			continue
		}
		addr, ok := d.lookup(ring[i])
		if !ok {
			continue
		}
		dest := h3.ToGeo(ring[i])
		neighbors = append(neighbors, Neighbor{
			Cell:      Cell{Host: addr, H3ID: ring[i]},
			DistanceM: h3.PointDistM(src, dest),
		})
	}
	sort.Slice(neighbors, func(i, j int) bool {
		return neighbors[i].DistanceM < neighbors[j].DistanceM
	})
	return
}

// ReplicaFor returns a list of hosts for replication.
func (d *Distributed) ReplicaFor(cell h3.H3Index, n int) ([]string, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if n > len(d.nodes) {
		return nil, fmt.Errorf("h3geodist: insufficient number of nodes want %d, have %d",
			n, len(d.nodes))
	}

	var mykey uint64
	var next int
	myaddr, ok := d.lookup(cell)
	if !ok {
		return nil, ErrVNodes
	}
	keys := make([]uint64, 0, 4)
	hosts := make(map[uint64]*node)
	for i := 0; i < len(d.nodes); i++ {
		hk := str2hash(d.nodes[i].addr)
		if d.nodes[i].addr == myaddr {
			mykey = hk
		}
		hosts[hk] = d.nodes[i]
		keys = append(keys, hk)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	res := make([]string, 0, 4)
	for next < len(keys) {
		if keys[next] == mykey {
			res = append(res, hosts[keys[next]].addr)
			break
		}
		next++
	}
	for len(res) < n {
		next++
		if next >= len(keys) {
			next = 0
		}
		res = append(res, hosts[keys[next]].addr)
	}
	return res, nil
}

// LookupMany returns a list of distributed cell.
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

// VNodeIndex returns the Index of the virtual node by H3Index.
func (d *Distributed) VNodeIndex(cell h3.H3Index) int {
	hashKey := uint2hash(uint64(cell))
	return int(hashKey % d.vnodes)
}

// EachVNode iterate each vnode, calling fn for each vnode.
func (d *Distributed) EachVNode(fn func(vnode uint64, addr string) bool) {
	for i := uint64(0); i < d.vnodes; i++ {
		addr, ok := d.Addr(i)
		if !ok {
			continue
		}
		if !fn(i, addr) {
			break
		}
	}
}

// Addr returns the addr of the node by vnode id.
func (d *Distributed) Addr(vnode uint64) (addr string, ok bool) {
	hashKey := uint2hash(vnode)
	idx := int(hashKey % d.vnodes)
	d.mu.RLock()
	node, found := d.index[idx]
	d.mu.RUnlock()
	if !found {
		return
	}
	addr = node.addr
	ok = true
	return
}

// EachCell iterate each distributed cell, calling fn for each cell.
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

// Add adds a new node.
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

// Remove removes a node.
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

// AvgLoad returns the average load.
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
