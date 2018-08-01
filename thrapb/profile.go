package thrapb

func DefaultProfile() *Profile {
	return &Profile{
		ID:           "default",
		Orchestrator: "docker",
	}
}
