package thrapb

import digest "github.com/opencontainers/go-digest"

// NewArtifact returns a Artifact instance with the id and tags
func NewArtifact(id string, tags []string) *Artifact {
	c := &Artifact{ID: digest.Digest(id)}
	c.SetTags(tags)
	return c
}

// SetTags sets the given tags
func (ci *Artifact) SetTags(tags []string) {
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
func (ci *Artifact) IsTagged() bool {
	return len(ci.Tags) > 0 && ci.Tags[0] != "<none>"
}
