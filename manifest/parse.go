package manifest

import (
	"io/ioutil"

	pb "github.com/euforia/thrap/thrapb"
	"github.com/hashicorp/hcl"
	"gopkg.in/yaml.v2"
)

// ParseYAML parses a manifest yaml file
func ParseYAML(file string) (mf *pb.Stack, err error) {
	var in []byte
	in, err = ioutil.ReadFile(file)
	if err == nil {
		mf, err = ParseYAMLBytes(in)
	}
	return
}

// ParseYAMLBytes reads a stack configuration into a stack struct
func ParseYAMLBytes(in []byte) (m *pb.Stack, err error) {
	var stack pb.Stack
	err = yaml.Unmarshal(in, &stack)
	if err == nil {
		m = &stack
	}
	return
}

// ParseHCL parses a hcl manifest file
func ParseHCL(file string) (mf *pb.Stack, err error) {
	var in []byte
	in, err = ioutil.ReadFile(file)
	if err == nil {
		mf, err = ParseHCLBytes(in)
	}
	return
}

func ParseHCLBytes(in []byte) (*pb.Stack, error) {
	// wrapper map
	var ms map[string]map[string]*pb.Stack
	err := hcl.Unmarshal(in, &ms)
	if err != nil {
		return nil, err
	}

	var stack *pb.Stack
	for _, v := range ms["manifest"] {
		stack = v
		break
	}
	return stack, nil
}
