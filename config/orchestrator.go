package config

// OrchestratorConfig holds oconfigurations for a orchestration
type OrchestratorConfig struct {
	ID   string `hcl:"id" hcle:"omit"`
	Addr string `hcl:"addr" hcle:"omitempty"`
}

// Clone returns a copy of the config
func (conf *OrchestratorConfig) Clone() *OrchestratorConfig {
	if conf == nil {
		return nil
	}
	return &OrchestratorConfig{
		ID:   conf.ID,
		Addr: conf.Addr,
	}
}

// Merge merges the other config into the one. Only non-empty fields are
// considered
func (conf *OrchestratorConfig) Merge(other *OrchestratorConfig) {
	if other == nil {
		return
	}

	if other.ID != "" {
		conf.ID = other.ID
	}

	if other.Addr != "" {
		conf.Addr = other.Addr
	}
}
