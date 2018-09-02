package thrapb

import (
	"encoding/binary"
	"errors"
	"hash"
	"math/rand"
	"time"

	"github.com/opencontainers/go-digest"

	"github.com/euforia/kvdb"
)

const (
	// ZeroSHA256Digest is a zero digest used on creation
	ZeroSHA256Digest = digest.Digest("sha256:0000000000000000000000000000000000000000000000000000000000000000")
)

// NewDeployment returns a new deployment object
func NewDeployment(profID, name string) *Deployment {
	return &Deployment{
		Name:      name,
		CreatedAt: time.Now().UnixNano(),
		Profile:   &Profile{ID: profID},
		Nonce:     rand.Uint64(),
		Version:   1,
		Previous:  ZeroSHA256Digest,
	}
}

// Clone retunrs a deep copy of a Deployment
func (d *Deployment) Clone() *Deployment {
	if d == nil {
		return nil
	}

	n := &Deployment{
		Previous:     d.Previous,
		Name:         d.Name,
		Version:      d.Version,
		CreatedAt:    d.CreatedAt,
		ModifiedAt:   d.ModifiedAt,
		Nonce:        d.Nonce,
		State:        d.State,
		Status:       d.Status,
		StateMessage: d.StateMessage,
	}
	if d.Spec != nil {
		n.Spec = make([]byte, len(d.Spec))
		copy(n.Spec, d.Spec)
	}
	n.Profile = d.Profile.Clone()

	return n
}

func (d *Deployment) PreviousDigest() digest.Digest {
	return d.Previous
}

// Hash satisfies the ObjectVersion interface
func (d *Deployment) Hash(h hash.Hash) {
	h.Write([]byte(d.Previous))
	h.Write([]byte(d.Name))
	binary.Write(h, binary.BigEndian, d.Version)
	binary.Write(h, binary.BigEndian, d.State)
	binary.Write(h, binary.BigEndian, d.Status)
	binary.Write(h, binary.BigEndian, d.CreatedAt)
	binary.Write(h, binary.BigEndian, d.ModifiedAt)
	binary.Write(h, binary.BigEndian, d.Nonce)
	h.Write(d.Spec)
	h.Write([]byte(d.StateMessage))

	d.Profile.Hash(h)
}

// New satisfies the ObjectVersion interface
func (d *Deployment) New() kvdb.ObjectVersion {
	return &Deployment{}
}

// Validate validates the deployment settings
func (d *Deployment) Validate() error {
	if d.Name == "" {
		return errors.New("name required")
	}
	if d.Profile == nil || len(d.Profile.ID) == 0 {
		return errors.New("profile id required")
	}
	return nil
}

// New satisfies the Object interface
func (desc *DeploymentDescriptor) New() kvdb.Object {
	return &DeploymentDescriptor{}
}
