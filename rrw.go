package h3geodist

type rrw struct {
	nodes []*node
	gcd   int
	maxw  int
	index int
	cw    int
	cnt   int
}

func newrrw() *rrw {
	return &rrw{nodes: make([]*node, 0, 8)}
}

func (r *rrw) exist(addr string) (ok bool) {
	for i := 0; i < len(r.nodes); i++ {
		if addr == r.nodes[i].addr {
			ok = true
			break
		}
	}
	return
}

func (r *rrw) addrAt(index int) string {
	return r.nodes[index].addr
}

func (r *rrw) remove(addr string) (ok bool) {
	for i := 0; i < len(r.nodes); i++ {
		if addr == r.nodes[i].addr {
			ok = true
			r.nodes = append(r.nodes[:i], r.nodes[i+1:]...)
			if r.cnt > 0 {
				r.cnt--
			}
			break
		}
	}
	if len(r.nodes) == 0 {
		r.nodes = r.nodes[:0]
		r.cnt = 0
		r.gcd = 0
		r.maxw = 0
		r.index = -1
		r.cw = 0
	}
	return
}

func (r *rrw) add(addr string, weight int) (ok bool) {
	if weight <= 0 {
		weight = 1
	}
	if r.exist(addr) {
		return
	}
	n := &node{addr: addr, weight: weight}
	if r.gcd > 0 {
		r.gcd = gcd(r.gcd, weight)
		if r.maxw < weight {
			r.maxw = weight
		}
	} else {
		r.gcd = weight
		r.maxw = weight
		r.index = -1
		r.cw = 0
	}
	r.nodes = append(r.nodes, n)
	r.cnt++
	ok = true
	return
}

func (r *rrw) size() int {
	return r.cnt
}

func (r *rrw) reset() {
	r.cw = 0
	r.index = -1
}

func (r *rrw) next() *node {
	if r.cnt == 0 {
		return nil
	}
	if r.cnt == 1 {
		return r.nodes[0]
	}
	for {
		r.index = (r.index + 1) % r.cnt
		if r.index == 0 {
			r.cw = r.cw - r.gcd
			if r.cw <= 0 {
				r.cw = r.maxw
				if r.cw == 0 {
					return nil
				}
			}
		}
		weight := r.nodes[r.index].weight
		if weight >= r.cw {
			return r.nodes[r.index]
		}
	}
}

func gcd(x, y int) int {
	var tmp int
	for {
		tmp = x % y
		if tmp > 0 {
			x, y = y, tmp
		} else {
			return y
		}
	}
}
