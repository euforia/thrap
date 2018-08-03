package thrapb

// DefaultProfile returns the default local profile
func DefaultProfile() *Profile {
	return &Profile{
		ID:           "local",
		Orchestrator: "docker",
		Registry:     "docker",
	}
}
