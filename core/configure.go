package core

import (
	"crypto/elliptic"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/euforia/hclencoder"
	"github.com/euforia/thrap/config"
	"github.com/euforia/thrap/consts"
	"github.com/euforia/thrap/utils"
	"github.com/euforia/thrap/vcs"
	homedir "github.com/mitchellh/go-homedir"
)

// ConfigureOptions holds options to configure a data directory
type ConfigureOptions struct {
	NoPrompt bool   // do not prompt for any input
	DataDir  string // data directory
	VCS      *config.VCSConfig
}

// DefaultConfigureOptions returns options with github as the vcs
func DefaultConfigureOptions() ConfigureOptions {
	return ConfigureOptions{
		VCS: &config.VCSConfig{ID: "github"},
	}
}

// ConfigureGlobal configures the global configuration
func ConfigureGlobal(opts ConfigureOptions) error {
	var err error

	if opts.DataDir == "" {
		opts.DataDir, err = homedir.Dir()
	} else if strings.HasPrefix(opts.DataDir, "~") {
		opts.DataDir, err = homedir.Expand(opts.DataDir)
	} else if !filepath.IsAbs(opts.DataDir) {
		opts.DataDir, err = filepath.Abs(opts.DataDir)
	}

	if err != nil {
		return err
	}

	// hwdir := filepath.Join(opts.DataDir, consts.WorkDir)
	if !utils.FileExists(opts.DataDir) {
		os.MkdirAll(opts.DataDir, 0755)
	}

	var (
		conf     *config.ThrapConfig
		varsfile = filepath.Join(opts.DataDir, consts.ConfigFile)
	)

	if !utils.FileExists(varsfile) {
		conf = config.DefaultThrapConfig()
	} else {
		conf, err = config.ReadThrapConfig(varsfile)
		if err != nil {
			return err
		}
	}

	if opts.VCS.Username != "" {
		if _, ok := conf.VCS[opts.VCS.ID]; !ok {
			return fmt.Errorf("unknown vcs provider: '%s'", opts.VCS.ID)
		}
		conf.VCS[opts.VCS.ID].Username = opts.VCS.Username
	}
	configureHomeVars(conf.VCS[opts.VCS.ID], opts.NoPrompt)

	err = config.WriteThrapConfig(conf, varsfile)
	if err != nil {
		return err
	}
	fmt.Println("Config:", varsfile)

	// Creds
	var (
		cconf     *config.CredsConfig
		credsFile = filepath.Join(opts.DataDir, consts.CredsFile)
	)
	if !utils.FileExists(credsFile) {
		cconf = config.DefaultCredsConfig()
	} else {
		cconf, err = config.ReadCredsConfig(credsFile)
		if err != nil {
			return err
		}
	}
	configureVCSCreds(cconf, opts.VCS.ID, opts.NoPrompt)
	err = config.WriteCredsConfig(cconf, credsFile)
	if err != nil {
		return err
	}
	fmt.Println("Creds:", credsFile)

	// Key file
	keypath := filepath.Join(opts.DataDir, consts.KeyFile)
	fmt.Println("Keypair:", keypath)
	if !utils.FileExists(keypath) {
		_, err = utils.GenerateECDSAKeyPair(keypath, elliptic.P256())
	}

	return err
}

func configureHomeVars(conf *config.VCSConfig, noprompt bool) {
	ghvcs, _ := vcs.New(&vcs.Config{Provider: "git"})
	if conf.Username == "" {
		conf.Username = ghvcs.GlobalUser()
	}

	if conf.Username != "" || noprompt {
		return
	}

	// Prompt for vcs.username
	var vname string
	prompt := fmt.Sprintf("%s username: ", conf.ID)
	utils.PromptUntilNoError(prompt, os.Stdout, os.Stdin, func(input []byte) error {
		vname = string(input)
		if vname == "" {
			return fmt.Errorf("%s username required", conf.ID)
		}
		return nil
	})

	conf.Username = vname

}

func configureVCSCreds(conf *config.CredsConfig, vcsID string, noprompt bool) {

	token := conf.VCS[vcsID]["token"]
	if token != "" || noprompt {
		return
	}

	prompt := fmt.Sprintf("%s token: ", vcsID)
	utils.PromptUntilNoError(prompt, os.Stdout, os.Stdin, func(input []byte) error {
		vcsToken := string(input)
		if vcsToken == "" {
			return fmt.Errorf("%s access token required", vcsID)
		}

		conf.VCS[vcsID]["token"] = vcsToken
		return nil
	})
}

// ConfigureLocal configures the local project configuration
func ConfigureLocal(conf *config.ThrapConfig, opts ConfigureOptions) (*config.ThrapConfig, error) {
	apath, err := filepath.Abs(opts.DataDir)
	if err != nil {
		return nil, err
	}

	tdir := filepath.Join(apath, consts.WorkDir)
	os.MkdirAll(tdir, 0755)

	varsfile := filepath.Join(tdir, consts.ConfigFile)
	if utils.FileExists(varsfile) {
		return config.ReadThrapConfig(varsfile)
	}

	// Read global config
	// filename, _ := homedir.Expand("~/" + consts.WorkDir + "/" + consts.ConfigFile)
	// conf, err := config.ReadThrapConfig(filename)
	// if err != nil {
	// 	return nil, err
	// }

	// Add project settings to supplied global
	conf.VCS[opts.VCS.ID].Repo = &config.VCSRepoConfig{
		Name:  opts.VCS.Repo.Name,
		Owner: opts.VCS.Repo.Owner,
	}

	b, err := hclencoder.Encode(conf)
	if err == nil {
		err = ioutil.WriteFile(varsfile, b, 0644)
	}

	return conf, err
}

// func generateKeyPair(filename string) (*ecdsa.PrivateKey, error) {
// 	c := elliptic.P256()
// 	kp, err := ecdsa.GenerateKey(c, rand.Reader)
// 	if err != nil {
// 		return nil, err
// 	}

// 	priv, pub, err := encodeECDSA(kp, filename)
// 	if err == nil {
// 		err = writePem(priv, pub, filename)
// 	}
// 	return kp, err
// }

// func encodeECDSA(privateKey *ecdsa.PrivateKey, filename string) ([]byte, []byte, error) {

// 	x509Encoded, err := x509.MarshalECPrivateKey(privateKey)
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	// pemEncoded := pem.Encode(privH, &pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})
// 	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})

// 	x509EncodedPub, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	pemEncodedPub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509EncodedPub})

// 	return pemEncoded, pemEncodedPub, nil
// }

// func writePem(priv, pub []byte, filename string) error {
// 	err := ioutil.WriteFile(filename, priv, 0600)
// 	if err != nil {
// 		return err
// 	}
// 	return ioutil.WriteFile(filename+".pub", pub, 0600)
// }

// func decodeECDSA(filename string) (*ecdsa.PrivateKey, error) {
// 	priv, err := ioutil.ReadFile(filename)
// 	if err != nil {
// 		return nil, err
// 	}

// 	block, _ := pem.Decode(priv)
// 	x509Encoded := block.Bytes
// 	privateKey, err := x509.ParseECPrivateKey(x509Encoded)
// 	if err != nil {
// 		return nil, err
// 	}

// 	pub, err := ioutil.ReadFile(filename + ".pub")
// 	if err != nil {
// 		return nil, err
// 	}
// 	blockPub, _ := pem.Decode(pub)
// 	// blockPub, _ := pem.Decode(pemEncodedPub)
// 	x509EncodedPub := blockPub.Bytes
// 	genericPublicKey, err := x509.ParsePKIXPublicKey(x509EncodedPub)
// 	if err != nil {
// 		return nil, err
// 	}

// 	publicKey := genericPublicKey.(*ecdsa.PublicKey)
// 	privateKey.PublicKey = *publicKey

// 	return privateKey, nil
// }
