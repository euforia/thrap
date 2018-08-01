package core

import (
	"github.com/docker/docker/api/types"
)

// CompStatus holds the overall component status
type CompStatus struct {
	ID      string
	Details types.ContainerJSON
	Error   error
}
