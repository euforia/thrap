package pb

import (
	"errors"
	"hash"
)

// DefaultProfile returns the default local profile
func DefaultProfile() *Profile {
	return &Profile{
		ID:           "local",
		Name:         "Local",
		Orchestrator: "docker",
		Registry:     "docker",
	}
}

// Validate validates the profile checking for required fields
func (p *Profile) Validate() error {
	if p.ID == "" {
		return errors.New("profile ID required")
	}
	if p.Name == "" {
		p.Name = p.ID
	}

	if p.Orchestrator == "" {
		return errors.New("profile orchestrator provider ID required")
	}
	if p.Secrets == "" {
		return errors.New("profile secrets provider ID required")
	}

	// We don't validate the registry as the default docker hub registry is
	// loaded and assumed
	return nil
}

// Hash writes the data structure contents to the hash function
func (p *Profile) Hash(h hash.Hash) {
	h.Write([]byte(p.ID))
	h.Write([]byte(p.Name))
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
		Name:         p.Name,
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
