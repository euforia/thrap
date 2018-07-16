package utils

import (
	"io/ioutil"

	"github.com/euforia/hclencoder"

	"github.com/euforia/thrap/thrapb"
	"github.com/hashicorp/hcl"
)

func LoadIdentities(filename string) (map[string]*thrapb.Identity, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var out map[string]*thrapb.Identity
	err = hcl.Unmarshal(b, &out)

	return out, err
}

func WriteIdentities(filename string, idents map[string]*thrapb.Identity) error {
	b, err := hclencoder.Encode(idents)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, b, 0644)
}
