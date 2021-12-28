package h3geodist

const (
	DefaultReplicationFactor = 9
	DefaultLoadFactor        = 1.25
	DefaultVNodes            = 256
	DefaultLevel             = 5
)

type Option func(*Distributed)

func WithVNodes(val uint64) Option {
	return func(d *Distributed) {
		d.vnodes = val
	}
}

func WithLoadFactor(val float64) Option {
	return func(d *Distributed) {
		d.loadFactor = val
	}
}

func WithReplicationFactor(val int) Option {
	return func(d *Distributed) {
		if val < 3 {
			val = 3
		}
		d.replFactor = val
	}
}
