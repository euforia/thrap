package config

// RepoConfig can be a registry repo or a VCS repo config
type RepoConfig struct {
	Name  string `hcl:"name"`
	Owner string `hcl:"owner"`
}

// Clone returns a copy of the config
func (conf *RepoConfig) Clone() *RepoConfig {
	if conf == nil {
		return nil
	}

	return &RepoConfig{
		Name:  conf.Name,
		Owner: conf.Owner,
	}
}

// Merge merges the other config into the one. Only non-empty fields are
// considered
func (conf *RepoConfig) Merge(other *RepoConfig) {
	if other == nil {
		return
	}

	if other.Name != "" {
		conf.Name = other.Name
	}
	if other.Owner != "" {
		conf.Owner = other.Owner
	}

}
