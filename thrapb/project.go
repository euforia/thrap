package thrapb

import (
	"hash"

	"github.com/euforia/kvdb"
)

// New satisfies the kvdb.Object interface
func (p *Project) New() kvdb.Object {
	return &Project{}
}

// Hash returns the hash of the project data
func (p *Project) Hash(h hash.Hash) []byte {
	h.Write([]byte(p.ID))
	h.Write([]byte(p.Name))
	h.Write([]byte(p.Description))
	h.Write([]byte(p.Source))
	h.Write([]byte(p.Owner))
	h.Write([]byte(p.Maintainer))
	for _, dev := range p.Developers {
		h.Write([]byte(dev))
	}
	return h.Sum(nil)
}

// Clone returns a copy of the project
func (p *Project) Clone() *Project {
	if p == nil {
		return p
	}
	proj := &Project{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Source:      p.Source,
		Owner:       p.Owner,
		Maintainer:  p.Maintainer,
		Developers:  make([]string, len(p.Developers)),
	}
	copy(proj.Developers, p.Developers)
	return proj
}
