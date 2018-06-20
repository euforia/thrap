package orchestrator

import "fmt"

// Config holds the config used to init the orchestrator
type Config struct {
	Provider string
	Conf     map[string]interface{}
}

// Orchestrator implements an application/project deployment orchestrator
type Orchestrator interface {
	Init(map[string]interface{}) error
	// NewSpec returns job spec/s based on the orchestrator.  There may be more
	// than one spec. e.g. kubernetes
	NewSpec(string) ([]string, error)
}

// New returns a new orchestrator based on the given config
func New(conf *Config) (Orchestrator, error) {
	var (
		orch Orchestrator
		err  error
	)

	switch conf.Provider {
	case "nomad":
		orch = &nomadOrchestrator{}

	default:
		err = fmt.Errorf("unsupported orchestrator: '%s'", conf.Provider)

	}

	if err == nil {
		err = orch.Init(conf.Conf)
	}

	return orch, err
}
