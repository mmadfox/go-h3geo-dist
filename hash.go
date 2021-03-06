package h3geodist

import (
	"encoding/binary"
	"hash/fnv"
)

// ToHash returns the fvn.Hash64 hash sum from uint64 value.
func ToHash(val uint64) uint64 {
	return uint2hash(val)
}

func str2hash(val string) uint64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(val))
	return h.Sum64()
}

func uint2hash(val uint64) uint64 {
	h := fnv.New64a()
	b64 := make([]byte, 8)
	binary.LittleEndian.PutUint64(b64, val)
	_, _ = h.Write(b64)
	return h.Sum64()
}
