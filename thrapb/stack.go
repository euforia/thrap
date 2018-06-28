package thrapb

import (
	"encoding/binary"
	"errors"
	"hash"
	"sort"

	"github.com/euforia/pseudo/scope"
	"github.com/euforia/thrap/vars"
	"github.com/hashicorp/hil/ast"
)

var (
	errDepCannotBuild = errors.New("dependencies cannot be built")
)

// Hash returns the hash of the header using the given hash function
func (header *ChainHeader) Hash(h hash.Hash) []byte {
	h.Reset()
	binary.Write(h, binary.BigEndian, header.Height)
	h.Write(header.Previous)
	binary.Write(h, binary.BigEndian, header.Timestamp)
	h.Write(header.DataDigest)
	return h.Sum(nil)
}

// ScopeVars returns the scoped variables usable for interpolation
func (stack *Stack) ScopeVars() scope.Variables {
	svars := scope.Variables{
		vars.StackDescription: ast.Variable{
			Value: stack.Description,
			Type:  ast.TypeString,
		},
		vars.StackName: ast.Variable{
			Value: stack.Name,
			Type:  ast.TypeString,
		},
		vars.StackID: ast.Variable{
			Value: stack.ID,
			Type:  ast.TypeString,
		},
		vars.StackVersion: ast.Variable{
			Value: stack.Version,
			Type:  ast.TypeString,
		},
	}

	for _, c := range stack.Components {
		sv := c.ScopeVars("comps.")
		for k, v := range sv {
			svars[k] = v
		}
	}
	for _, c := range stack.Dependencies {
		sv := c.ScopeVars("deps.")
		for k, v := range sv {
			svars[k] = v
		}
	}

	return svars
}

// Hash returns the hash of the object using the given hash function
func (stack *Stack) Hash(h hash.Hash) []byte {

	h.Write([]byte(stack.ID))
	h.Write([]byte(stack.Name))
	h.Write([]byte(stack.Version))
	h.Write([]byte(stack.Description))

	keys := make([]string, 0, len(stack.Components))
	for k := range stack.Components {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		h.Write([]byte(k))
		stack.Components[k].Hash(h)
	}

	keys = make([]string, 0, len(stack.Dependencies))
	for k := range stack.Dependencies {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		h.Write([]byte(k))
		stack.Dependencies[k].Hash(h)

	}

	return h.Sum(nil)
}

// Validate validates the stack configuration.  On success all component and
// dependency ids are set otherwise an error is returned
func (stack *Stack) Validate() map[string]error {

	if stack.ID == "" {
		stack.ID = stack.Name
	}

	errs := make(map[string]error)

	for k, comp := range stack.Components {
		if err := comp.Validate(); err != nil {
			errs["component."+k] = err
		} else {
			comp.ID = k
		}
	}

	for k, dep := range stack.Dependencies {
		if err := validateDep(dep); err != nil {
			errs["dependency."+k] = err
		} else {
			dep.ID = k
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil

}

func validateDep(dep *Component) error {
	if dep.IsBuildable() {
		return errDepCannotBuild
	}
	return nil
}
