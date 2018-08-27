package thrapb

import (
	"encoding/binary"
	"hash"
)

// Hash writes the data structure contents to the hash function
func (hc *HealthCheck) Hash(h hash.Hash) {
	h.Write([]byte(hc.Protocol))
	h.Write([]byte(hc.Path))
	h.Write([]byte(hc.Method))
	binary.Write(h, binary.BigEndian, hc.Timeout)
	binary.Write(h, binary.BigEndian, hc.Interval)
	h.Write([]byte(hc.PortLabel))
}
