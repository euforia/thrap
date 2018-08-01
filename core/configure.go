package core

import (
	"crypto/elliptic"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/euforia/hclencoder"
	"github.com/euforia/thrap/config"
	"github.com/euforia/thrap/consts"
	"github.com/euforia/thrap/utils"
	"github.com/euforia/thrap/vcs"
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
	opts.DataDir, err = utils.GetAbsPath(opts.DataDir)
	if err != nil {
		return err
	}

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
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.StripEscape)
	fmt.Fprintf(tw, "Config:\t%s\n", varsfile)

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
	fmt.Fprintf(tw, "Creds:\t%s\n", credsFile)

	// Key file
	keypath := filepath.Join(opts.DataDir, consts.KeyFile)
	if !utils.FileExists(keypath) {
		_, err = utils.GenerateECDSAKeyPair(keypath, elliptic.P256())
	}
	fmt.Fprintf(tw, "Keypair:\t%s\n", keypath)
	tw.Flush()

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
