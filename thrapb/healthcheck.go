package thrapb

import (
	"encoding/binary"
	"hash"
)

func (hc *HealthCheck) Hash(h hash.Hash) {
	h.Write([]byte(hc.Protocol))
	h.Write([]byte(hc.Path))
	h.Write([]byte(hc.Method))
	binary.Write(h, binary.BigEndian, hc.Timeout)
	binary.Write(h, binary.BigEndian, hc.Interval)
	h.Write([]byte(hc.PortLabel))
}
