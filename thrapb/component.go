package thrapb

import (
	"encoding/binary"
	"errors"
	"hash"
	"sort"

	"github.com/euforia/pseudo/scope"
	"github.com/euforia/thrap/consts"
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

// NewComponent returns a new Compoenent of the given type, name and version
func NewComponent(name, version string, typ CompType) *Component {
	return &Component{
		Name:    name,
		Type:    typ,
		Version: version,
	}
}

// ScopeVarName returns the scoped variable name for this component for the
// given key
func (comp *Component) ScopeVarName(prefix, key string) string {
	return prefix + comp.ID + "." + key
}

// ScopeVars returns the scoped variables using the prefix for the key
func (comp *Component) ScopeVars(prefix string) scope.Variables {
	// prefix + comp.ID + ".version"
	svars := scope.Variables{
		comp.ScopeVarName(prefix, "version"): ast.Variable{
			Type:  ast.TypeString,
			Value: comp.Version,
		},
	}

	for k, p := range comp.Ports {
		vname := comp.ScopeVarName(prefix, "container.port."+k)
		svars[vname] = ast.Variable{
			Value: int(p),
			Type:  ast.TypeInt,
		}
	}

	return svars
}

// HasPort returns true if the component has the given port specified
func (comp *Component) HasPort(port int32) bool {
	for _, p := range comp.Ports {
		if p == port {
			return true
		}
	}

	return false
}

func (comp *Component) HasVolumeTarget(target string) bool {
	for _, vol := range comp.Volumes {
		if vol.Target == target {
			return true
		}
	}
	return false
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

// HasEnvVars returns true if the component has environment variables
// defined
func (comp *Component) HasEnvVars() bool {
	return comp.Env != nil && len(comp.Env.Vars) > 0
}

// HasLanguage returns true it the component has a programming language
// defined
func (comp *Component) HasLanguage() bool {
	return comp.Language.Lang() != ""
}

func (comp *Component) validateCommon() error {

	var err error

	if comp.IsBuildable() {
		// Make sure the language is valid
		if comp.HasLanguage() {
			err = comp.Language.Validate()
			if err != nil {
				return err
			}
		}

		if comp.Build.Context == "" {
			comp.Build.Context = consts.DefaultBuildContext
		}

	} else {
		// Check component version is provided if not a buildable
		_, err = version.NewVersion(comp.Version)
	}

	return err
}

// Hash computes the hash of the component to an existing hash writer. It
// does not reset the hasher or perform the final Sum(nil)
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
		h.Write([]byte(comp.Secrets.Template))
	}

	if comp.Ports != nil {
		keys := make([]string, 0, len(comp.Ports))
		for k := range comp.Ports {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			h.Write([]byte(k))
			binary.Write(h, binary.BigEndian, comp.Ports[k])
		}
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

	for _, vol := range comp.Volumes {
		h.Write([]byte(vol.Source))
		h.Write([]byte(vol.Target))
	}
}
