package config

// SecretsConfig holds secrets provider configurations
type SecretsConfig struct {
	ID   string `hcl:"id"   hcle:"omit"`
	Addr string `hcl:"addr" hcle:"omitempty"`
}

// Clone returns a copy of the config
func (conf *SecretsConfig) Clone() *SecretsConfig {
	if conf == nil {
		return nil
	}
	return &SecretsConfig{ID: conf.ID, Addr: conf.Addr}
}

// Merge merges the other config into the one. Only non-empty fields are
// considered
func (conf *SecretsConfig) Merge(other *SecretsConfig) {
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
