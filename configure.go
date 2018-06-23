package thrap

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/euforia/hclencoder"
	"github.com/euforia/thrap/config"
	"github.com/euforia/thrap/consts"
	"github.com/euforia/thrap/utils"
	"github.com/euforia/thrap/vcs"
	homedir "github.com/mitchellh/go-homedir"
)

// ConfigureHomeDir checks for the global working directory
func ConfigureHomeDir(noprompt bool) error {
	hdir, err := homedir.Dir()
	if err != nil {
		return err
	}

	hwdir := filepath.Join(hdir, consts.WorkDir)
	if !utils.FileExists(hwdir) {
		os.MkdirAll(hwdir, 0755)
	}

	// conf := config.DefaultThrapConfig()
	var conf *config.ThrapConfig

	varsfile := filepath.Join(hwdir, consts.ConfigFile)
	if !utils.FileExists(varsfile) {
		conf = config.DefaultThrapConfig()
	} else {
		conf, err = ReadGlobalConfig()
		if err != nil {
			return err
		}
	}
	// fmt.Println(conf.VCS, conf.VCS["github"])
	configureHomeVars(conf.VCS["github"], varsfile, noprompt)

	b, err := hclencoder.Encode(conf)
	if err == nil {
		err = ioutil.WriteFile(varsfile, b, 0644)
	}

	var cconf *config.CredsConfig

	credsFile := filepath.Join(hwdir, consts.CredsFile)
	if !utils.FileExists(credsFile) {
		cconf = config.DefaultCredsConfig()
	} else {
		cconf, err = ReadGlobalCreds()
		if err != nil {
			return err
		}
	}

	err = configureVCSCreds(cconf, "github", credsFile)
	if err != nil {
		return err
	}

	fpath, _ := homedir.Expand("~/" + consts.WorkDir + "/" + consts.KeyFile)
	if !utils.FileExists(fpath) {
		_, err = generateKeyPair(fpath)
	}

	return err
}

func configureHomeVars(conf *config.VCSConfig, varsfile string, noprompt bool) {
	ghvcs, _ := vcs.New(&vcs.Config{Provider: "git"})
	if conf.Username == "" {
		conf.Username = ghvcs.GlobalUser()
	}

	if conf.Username != "" || noprompt {
		return
	}

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

func configureVCSCreds(conf *config.CredsConfig, vcsID, credsFile string) error {

	token := conf.VCS[vcsID]["token"]
	if token == "" {

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

	return config.WriteCredsConfig(conf, credsFile)
}

// ConfigureProjectDir sets up the project .thrap dir. name is the repo name
func ConfigureProjectDir(name, vcsp, repoOwner, dir string) (*config.ThrapConfig, error) {
	apath, err := filepath.Abs(dir)
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
	filename, _ := homedir.Expand("~/" + consts.WorkDir + "/" + consts.ConfigFile)
	conf, err := config.ReadThrapConfig(filename)
	if err != nil {
		return nil, err
	}

	// Add project settings
	conf.VCS[vcsp].Repo = &config.VCSRepoConfig{
		Name:  name,
		Owner: repoOwner,
	}

	b, err := hclencoder.Encode(conf)
	if err == nil {
		err = ioutil.WriteFile(varsfile, b, 0644)
	}

	return conf, err
}

func generateKeyPair(filename string) (*ecdsa.PrivateKey, error) {
	c := elliptic.P256()
	kp, err := ecdsa.GenerateKey(c, rand.Reader)
	if err != nil {
		return nil, err
	}

	priv, pub, err := encodeECDSA(kp, filename)
	if err == nil {
		err = writePem(priv, pub, filename)
	}
	return kp, err
}

func encodeECDSA(privateKey *ecdsa.PrivateKey, filename string) ([]byte, []byte, error) {

	x509Encoded, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return nil, nil, err
	}
	// pemEncoded := pem.Encode(privH, &pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})
	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})

	x509EncodedPub, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, nil, err
	}

	pemEncodedPub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509EncodedPub})

	return pemEncoded, pemEncodedPub, nil
}

func writePem(priv, pub []byte, filename string) error {
	err := ioutil.WriteFile(filename, priv, 0600)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename+".pub", pub, 0600)
}

func decodeECDSA(filename string) (*ecdsa.PrivateKey, error) {
	priv, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(priv)
	x509Encoded := block.Bytes
	privateKey, err := x509.ParseECPrivateKey(x509Encoded)
	if err != nil {
		return nil, err
	}

	pub, err := ioutil.ReadFile(filename + ".pub")
	if err != nil {
		return nil, err
	}
	blockPub, _ := pem.Decode(pub)
	// blockPub, _ := pem.Decode(pemEncodedPub)
	x509EncodedPub := blockPub.Bytes
	genericPublicKey, err := x509.ParsePKIXPublicKey(x509EncodedPub)
	if err != nil {
		return nil, err
	}

	publicKey := genericPublicKey.(*ecdsa.PublicKey)
	privateKey.PublicKey = *publicKey

	return privateKey, nil
}
