package h3geodist

const (
	DefaultReplicationFactor = 9
	DefaultLoadFactor        = 1.25
	DefaultVNodes            = 64
)

// Option is a type to represent various Distributed options.
type Option func(*Distributed)

// WithVNodes sets the number of virtual nodes. Default 64.
func WithVNodes(val uint64) Option {
	return func(d *Distributed) {
		d.vnodes = val
	}
}

// WithLoadFactor sets the number of load factor. Default 1.25.
func WithLoadFactor(val float64) Option {
	return func(d *Distributed) {
		d.loadFactor = val
	}
}

// WithReplicationFactor sets the number of replication factor. Default 9.
// A value less than or equal to 0 is set to 1.
func WithReplicationFactor(val int) Option {
	return func(d *Distributed) {
		if val <= 0 {
			val = 1
		}
		d.replFactor = val
	}
}
