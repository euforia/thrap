package thrapb

import "hash"

// DefaultProfile returns the default local profile
func DefaultProfile() *Profile {
	return &Profile{
		ID:           "local",
		Orchestrator: "docker",
		Registry:     "docker",
	}
}

func (p *Profile) Hash(h hash.Hash) {
	h.Write([]byte(p.ID))
	h.Write([]byte(p.Orchestrator))
	h.Write([]byte(p.Secrets))
	h.Write([]byte(p.Registry))
	h.Write([]byte(p.VCS))

	keys := make([]string, 0, len(p.Meta))
	for k := range p.Meta {
		keys = append(keys, k)
	}

	for _, k := range keys {
		h.Write([]byte(k))
		h.Write([]byte(p.Meta[k]))
	}

	keys = make([]string, 0, len(p.Variables))
	for k := range p.Variables {
		keys = append(keys, k)
	}

	for _, k := range keys {
		h.Write([]byte(k))
		h.Write([]byte(p.Variables[k]))
	}
}

// Clone returns a deep copy of the Profile
func (p *Profile) Clone() *Profile {
	if p == nil {
		return p
	}
	n := &Profile{
		ID:           p.ID,
		Orchestrator: p.Orchestrator,
		Registry:     p.Registry,
		Secrets:      p.Secrets,
		VCS:          p.VCS,
	}

	if p.Meta != nil {
		n.Meta = make(map[string]string)
		for k, v := range p.Meta {
			n.Meta[k] = v
		}
	}

	if p.Variables != nil {
		n.Variables = make(map[string]string)
		for k, v := range p.Variables {
			n.Variables[k] = v
		}
	}

	return n
}
