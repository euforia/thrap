package thrap

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"io/ioutil"
	"math/big"
	"path/filepath"

	"github.com/euforia/pseudo/scope"
	"github.com/euforia/thrap/config"
	"github.com/euforia/thrap/consts"
	homedir "github.com/mitchellh/go-homedir"
)

// func LoadManifest() (*thrapb.Stack, error) {
// 	mfile, err := utils.GetLocalPath(consts.DefaultManifestFile)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return manifest.LoadManifest(mfile)
// }

func ReadGlobalConfig() (*config.ThrapConfig, error) {
	filename, err := homedir.Expand("~/" + consts.WorkDir + "/" + consts.ConfigFile)
	if err == nil {
		return config.ReadThrapConfig(filename)
	}
	return nil, err
}

func ReadProjectConfig(projPath string) (*config.ThrapConfig, error) {
	filename := filepath.Join(projPath, consts.WorkDir, consts.ConfigFile)
	return config.ReadThrapConfig(filename)
}

func ReadGlobalCreds() (*config.CredsConfig, error) {
	filename, err := homedir.Expand("~/" + consts.WorkDir + "/" + consts.CredsFile)
	if err == nil {
		return config.ReadCredsConfig(filename)
	}
	return nil, err
}

func LoadGlobalScopeVars() (scope.Variables, error) {
	hdir, err := homedir.Dir()
	if err == nil {
		fpath := filepath.Join(hdir, consts.WorkDir, consts.ConfigFile)
		return LoadVariablesFromFile(fpath)
	}
	return nil, err
}

func LoadProjectScopeVars(projPath string) (scope.Variables, error) {
	return LoadVariablesFromFile(filepath.Join(projPath, consts.WorkDir, consts.ConfigFile))
}

func LoadVariablesFromFile(name string) (scope.Variables, error) {
	in, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	return LoadVariables(in, "")
}

func LoadVariables(in []byte, prefix string) (scope.Variables, error) {
	pbld := scope.NewPseudoBuilder(prefix)
	err := pbld.Build(in)
	if err == nil {
		return pbld.Variables(), nil
	}

	return nil, err

}

func LoadUserKeyPair() (*ecdsa.PrivateKey, error) {
	filename, err := homedir.Expand("~/" + consts.WorkDir + "/" + consts.KeyFile)
	if err == nil {
		return decodeECDSA(filename)
	}
	return nil, err
}

func verifySignature(pubkey, data, signature []byte) bool {
	r := big.Int{}
	s := big.Int{}
	sigLen := len(signature)
	r.SetBytes(signature[:(sigLen / 2)])
	s.SetBytes(signature[(sigLen / 2):])

	x := big.Int{}
	y := big.Int{}
	keyLen := len(pubkey)
	x.SetBytes(pubkey[:(keyLen / 2)])
	y.SetBytes(pubkey[(keyLen / 2):])

	rawPubKey := &ecdsa.PublicKey{Curve: elliptic.P256(), X: &x, Y: &y}

	return ecdsa.Verify(rawPubKey, data, &r, &s)
}

func makePubKeyFromBytes(curve elliptic.Curve, pubkey []byte) *ecdsa.PublicKey {
	x := big.Int{}
	y := big.Int{}
	keyLen := len(pubkey)
	x.SetBytes(pubkey[:(keyLen / 2)])
	y.SetBytes(pubkey[(keyLen / 2):])

	return &ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
}
