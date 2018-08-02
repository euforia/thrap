package thrapb

// DefaultProfile returns the default local profile
func DefaultProfile() *Profile {
	return &Profile{
		ID:           "default",
		Orchestrator: "docker",
		Registry:     "docker",
	}
}
