package core

import (
	"time"

	"github.com/docker/docker/api/types"
	digest "github.com/opencontainers/go-digest"
)

// CompImage holds image info about a component
type CompImage struct {
	ID      digest.Digest
	Tags    []string
	Labels  map[string]string
	Created time.Time
	Size    int64
}

// NewCompImage returns a CompImage instance with the id and tags
func NewCompImage(id string, tags []string) *CompImage {
	c := &CompImage{ID: digest.Digest(id)}
	c.SetTags(tags)
	return c
}

// SetTags sets the given tags
func (ci *CompImage) SetTags(tags []string) {
	ci.Tags = make([]string, 0, len(tags))
	for _, tag := range tags {
		if tag == "<none>:<none>" {
			ci.Tags = append(ci.Tags, "<none>")
		} else {
			ci.Tags = append(ci.Tags, tag)
		}
	}
}

// IsTagged returns true if the image has valid tags
func (ci *CompImage) IsTagged() bool {
	return len(ci.Tags) > 0 && ci.Tags[0] != "<none>"
}

// CompStatus holds the overall component status
type CompStatus struct {
	ID      string
	Details types.ContainerJSON
	Error   error
}
