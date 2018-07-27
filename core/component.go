package core

import (
	"github.com/docker/docker/api/types"
	"github.com/euforia/thrap/metrics"
)

// CompStatus holds the overall component status
type CompStatus struct {
	ID      string
	Details types.ContainerJSON
	Error   error
}

type CompBuildResult struct {
	ID      string
	Labels  map[string]string
	Runtime *metrics.Runtime
}
