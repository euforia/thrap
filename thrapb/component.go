package thrapb

import (
	"encoding/binary"
	"errors"
	"hash"
	"sort"

	"github.com/euforia/pseudo/scope"
	version "github.com/hashicorp/go-version"
	"github.com/hashicorp/hil/ast"
)

// CompType is the component type
type CompType string

const (
	// CompTypeUnknown is an unknown comp type
	CompTypeUnknown = ""
	// CompTypeWeb represents a ui component of the stack
	CompTypeWeb CompType = "web"
	// CompTypeAPI represents an api component of the stack
	CompTypeAPI CompType = "api"
	// CompTypeDatastore is a datastore component. This includes dataabases,
	// caches, storage and similar constructs
	CompTypeDatastore CompType = "datastore"
	// CompTypeBatch are components such as batch jobs, aws lamdba functions
	// and similar constructs
	CompTypeBatch CompType = "batch"
	// CompTypePeriodic are components like cron
	CompTypePeriodic CompType = "periodic"
)

var (
	errTypeNotSpecified = errors.New("component type not specified")
	errDatastoreHead    = errors.New("datastore cannot be a head")
)

func NewComponent(name, version string, typ CompType) *Component {
	return &Component{
		Name:    name,
		Type:    typ,
		Version: version,
	}
}

func (comp *Component) ScopeVars(prefix string) scope.Variables {

	return scope.Variables{
		prefix + "version": ast.Variable{
			Type:  ast.TypeString,
			Value: comp.Version,
		},
	}
}

// IsBuildable returns true if a build file has been specified signifying the
// component is buildable
func (comp *Component) IsBuildable() bool {
	return comp.Build != nil && comp.Build.Dockerfile != ""
}

// HasSecrets returns true if this component has specified secrets mgmt
func (comp *Component) HasSecrets() bool {
	return comp.Secrets != nil && comp.Secrets.Destination != ""
}

// Validate validates the component configuration
func (comp *Component) Validate() error {
	switch comp.Type {
	case CompTypeWeb, CompTypeAPI:

	case CompTypeDatastore:
		// We don't allow a datastore to be directly exposed
		if comp.Head {
			return errDatastoreHead
		}

	case CompTypeBatch, CompTypePeriodic:

	default:
		return errTypeNotSpecified

	}

	return comp.validateCommon()
}

func (comp *Component) validateCommon() error {

	var err error

	if comp.IsBuildable() {
		// Make sure the language is valid
		err = comp.Language.Validate()
	} else {
		// Check component version is provided if not a buildable
		_, err = version.NewVersion(comp.Version)
	}

	return err
}

func (comp *Component) Hash(h hash.Hash) {
	h.Write([]byte(comp.ID))
	h.Write([]byte(comp.Name))
	h.Write([]byte(comp.Version))
	h.Write([]byte(comp.Type))
	h.Write([]byte(comp.Language))

	binary.Write(h, binary.BigEndian, comp.Head)
	binary.Write(h, binary.BigEndian, comp.External)

	if comp.Build != nil {
		h.Write([]byte(comp.Build.Dockerfile))
		h.Write([]byte(comp.Build.Context))
	}
	if comp.Secrets != nil {
		h.Write([]byte(comp.Secrets.Destination))
		h.Write([]byte(comp.Secrets.Format))
	}
	if comp.Env != nil {
		h.Write([]byte(comp.Env.File))
		if comp.Env.Vars != nil {
			keys := make([]string, 0, len(comp.Env.Vars))
			for k := range comp.Env.Vars {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				h.Write([]byte(k + comp.Env.Vars[k]))
			}
		}
	}

	if comp.Config != nil {
		keys := make([]string, 0, len(comp.Config))
		for k := range comp.Config {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			h.Write([]byte(k + comp.Config[k]))
		}
	}
}
