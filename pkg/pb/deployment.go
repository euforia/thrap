package pb

import (
	"encoding/binary"
	"errors"
	"hash"
	"math/rand"
	"time"

	"github.com/opencontainers/go-digest"
)

const (
	// ZeroSHA256Digest is a zero digest used on creation
	ZeroSHA256Digest = digest.Digest("sha256:0000000000000000000000000000000000000000000000000000000000000000")
)

const (
	// DescContentTypeMoldHCL is the legacy type to be deprecated
	DescContentTypeMoldHCL = "application/vnd.thrap.mold.deployment.descriptor.v1+hcl"
	// DescContentTypeNomadHCL is a nomad hcl file
	DescContentTypeNomadHCL = "application/vnd.thrap.nomad.deployment.descriptor.v1+hcl"
	// DescContentTypeNomadJSON is json object
	DescContentTypeNomadJSON = "application/vnd.thrap.nomad.deployment.descriptor.v1+json"
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
		State:     DeploymentState_CREATE,
		Status:    DeployStateStatus_OK,
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
	n.Desc = d.Desc.Clone()
	n.Profile = d.Profile.Clone()

	return n
}

// func (d *Deployment) PreviousDigest() digest.Digest {
// 	return d.Previous
// }

// Hash satisfies the ObjectVersion interface
func (d *Deployment) Hash(h hash.Hash) {
	h.Write([]byte(d.Previous))
	h.Write([]byte(d.Name))
	if d.Desc != nil {
		h.Write([]byte(d.Desc.ID))
		h.Write([]byte(d.Desc.Mime))
		h.Write(d.Desc.Spec)
	}

	binary.Write(h, binary.BigEndian, d.Version)
	binary.Write(h, binary.BigEndian, d.State)
	binary.Write(h, binary.BigEndian, d.Status)
	binary.Write(h, binary.BigEndian, d.CreatedAt)
	binary.Write(h, binary.BigEndian, d.ModifiedAt)
	binary.Write(h, binary.BigEndian, d.Nonce)

	h.Write([]byte(d.StateMessage))

	d.Profile.Hash(h)
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

// Clone clones the DeploymentDescriptor
func (desc *DeploymentDescriptor) Clone() *DeploymentDescriptor {
	if desc == nil {
		return nil
	}
	n := &DeploymentDescriptor{
		ID:   desc.ID,
		Mime: desc.Mime,
	}
	if desc.Spec != nil {
		n.Spec = make([]byte, len(desc.Spec))
		copy(n.Spec, desc.Spec)
	}
	return n
}
