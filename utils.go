package thrap

import (
	"bufio"
	"crypto/ecdsa"
	"crypto/elliptic"
	"errors"
	"io"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"strings"

	"github.com/euforia/pseudo/scope"
	"github.com/euforia/thrap/config"
	"github.com/euforia/thrap/consts"
	"github.com/euforia/thrap/manifest"
	"github.com/euforia/thrap/thrapb"
	homedir "github.com/mitchellh/go-homedir"
)

func ReadGlobalConfig() (*config.ThrapConfig, error) {
	filename, err := homedir.Expand("~/" + consts.WorkDir + "/" + consts.ConfigFile)
	if err == nil {
		return config.ReadThrapConfig(filename)
	}
	return nil, err
}

func ReadGlobalCreds() (*config.CredsConfig, error) {
	filename, err := homedir.Expand("~/" + consts.WorkDir + "/" + consts.CredsFile)
	if err == nil {
		return config.ReadCredsConfig(filename)
	}
	return nil, err
}

func FlattenErrors(errs map[string]error) error {
	var out string
	for k, v := range errs {
		out += k + ":" + v.Error() + "\n"
	}
	return errors.New(out)
}

func FileExists(fpath string) bool {
	_, err := os.Stat(fpath)
	return err == nil
}

func mergeErrors(e1, e2 error) error {
	if e1 == nil {
		return e2
	} else if e2 == nil {
		return e1
	}
	return errors.New(e1.Error() + "; " + e2.Error())
}

// func LoadGlobalVars() (map[string]string, error) {
// 	hdir, err := homedir.Dir()
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return LoadVars(hdir)
// }
//
// func LoadVars(dir string) (map[string]string, error) {
// 	b, err := ioutil.ReadFile(filepath.Join(dir, consts.WorkDir, consts.DefaultVarsFile))
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	var out map[string]string
// 	err = hcl.Decode(&out, string(b))
// 	return out, err
// }
//
// func MergeVars(m1, m2 map[string]string) map[string]string {
// 	if m1 == nil {
// 		return m2
// 	} else if m2 == nil {
// 		return m1
// 	}
//
// 	for k, v := range m2 {
// 		m1[k] = v
// 	}
// 	return m1
// }

// GetLocalPath computes the path from the user specified args.  Uses the
// current directory if none is supplied in args
func GetLocalPath(in string) (dirpath string, err error) {
	// Assume cwd
	if len(in) == 0 {
		return os.Getwd()
	}

	// Assume cwd + supplied path if not an absolute path
	if !filepath.IsAbs(in) {
		var wd string
		if wd, err = os.Getwd(); err == nil {
			dirpath = filepath.Join(wd, in)
		}
	}

	return
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

func PromptUntilNoError(prompt string, out io.Writer, in io.Reader, f func([]byte) error) {
	var (
		lb []byte
	)
	err := io.ErrUnexpectedEOF
	for err != nil {
		out.Write([]byte(prompt))

		rd := bufio.NewReader(in)
		lb, _ = rd.ReadBytes('\n')
		lb = lb[:len(lb)-1]
		err = f(lb)
	}
}

// LoadManifest loads a thrap manifest.  A manifest begins with
// `manifest "id" {` followed by the remainder definition
func LoadManifest(mfile string) (*thrapb.Stack, error) {

	if mfile == "" {
		if FileExists(consts.DefaultManifestFile) {
			mfile = "thrap.hcl"
		} else if FileExists("thrap.yml") {
			mfile = "thrap.yml"
		} else {
			return nil, errors.New("no manifest found")
		}
	}

	mpath, err := GetLocalPath(mfile)
	if err != nil {
		return nil, err
	}

	var st *thrapb.Stack
	if strings.HasSuffix(mfile, ".hcl") {
		st, err = manifest.ParseHCL(mpath)
	} else {
		st, err = manifest.ParseYAML(mpath)
	}

	return st, err
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
