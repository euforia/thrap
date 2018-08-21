package thrapb

// DefaultProfile returns the default local profile
func DefaultProfile() *Profile {
	return &Profile{
		ID:           "local",
		Orchestrator: "docker",
		Registry:     "docker",
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

	if n.Meta != nil {
		n.Meta = make(map[string]string)
		for k, v := range p.Meta {
			n.Meta[k] = v
		}
	}

	if n.Variables != nil {
		n.Variables = make(map[string]string)
		for k, v := range p.Variables {
			n.Variables[k] = v
		}
	}

	return n
}
