package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/urfave/cli.v2"

	"github.com/euforia/thrap"
	"github.com/euforia/thrap/analysis"
	"github.com/euforia/thrap/asm"
	"github.com/euforia/thrap/consts"
	"github.com/euforia/thrap/packs"
	"github.com/euforia/thrap/thrapb"
	"github.com/euforia/thrap/utils"
	"github.com/euforia/thrap/vars"
	"github.com/euforia/thrap/vcs"
	"gopkg.in/src-d/go-git.v4"
)

var usageTextInit = `thrap init [command options] [directory]

   Init bootstraps a new project in the specified directory.  If no directory is
   given, it defaults to the current directory.

   It sets up the VCS, registries, secrets and any other configured resources.`

func commandStackInit() *cli.Command {
	return &cli.Command{
		Name:      "init",
		Usage:     "Initialize a new project",
		UsageText: usageTextInit,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "name",
				Aliases: []string{"n"},
				Usage:   "project `name` (default: <current directory>)",
			},
			&cli.StringFlag{
				Name:    "lang",
				Aliases: []string{"l"},
				Usage:   "programming `language`",
			},
			&cli.StringFlag{
				Name:  "vcs",
				Usage: "version control `provider`",
				Value: "github",
			},
			&cli.StringFlag{
				Name:  vars.VcsRepoOwner,
				Usage: "source code repo `owner`",
			},
		},
		Action: func(ctx *cli.Context) error {

			//
			projPath, err := setupProjPath(ctx)
			if err != nil {
				return err
			}

			// Project name
			projName := ctx.String("name")
			if len(projName) == 0 {
				projName = filepath.Base(projPath)
			}

			pks := packs.New("./etc/packs")
			devpacks := pks.Dev()

			// Set language from input or otherwise and other related params
			_, err = setLanguage(ctx, devpacks, projPath)
			if err != nil {
				return err
			}

			gconf, err := thrap.ReadGlobalConfig()
			if err != nil {
				return err
			}

			vcsID := ctx.String("vcs")
			defaultVCS := gconf.VCS[vcsID]
			if err = setRepoOwner(ctx, defaultVCS.Username, projPath); err != nil {
				return err
			}

			repoOwner := ctx.String(vars.VcsRepoOwner)

			// Local project setup
			pconf, err := thrap.ConfigureProjectDir(projName, defaultVCS.ID, repoOwner, projPath)
			if err != nil {
				return err
			}
			defaultVCS = pconf.VCS[vcsID]

			mfile := filepath.Join(projPath, consts.DefaultManifestFile)
			if utils.FileExists(mfile) {
				// TODO: ???
				return fmt.Errorf("manifest %s already exists", consts.DefaultManifestFile)
			}

			//gitRemoteAddr := scopeVars[vars.VcsAddr].Value.(string)
			vcsp, gitRepo, err := setupGit(projName, repoOwner, projPath, defaultVCS.Addr)
			if err != nil {
				return err
			}

			// conf, err := config.ParseFile("./etc/config.hcl")
			// if err != nil {
			// 	return err
			// }

			// Prompt for missing
			stack, err := promptComps(projName, ctx.String("lang"), pks)
			if err != nil {
				return err
			}
			fmt.Println()

			errs := stack.Validate()
			if errs != nil {
				return utils.FlattenErrors(errs)
			}

			scopeVars := defaultVCS.ScopeVars("vcs.")
			stAsm, err := asm.NewStackAsm(stack, vcsp, gitRepo, scopeVars, devpacks)
			if err != nil {
				return err
			}

			err = stAsm.Assemble()
			if err != nil {
				return err
			}

			return stAsm.WriteManifest()
		},
	}
}

func setupProjPath(ctx *cli.Context) (string, error) {
	var projPath string
	if args := ctx.Args(); args.Len() > 0 {
		projPath = args.First()
	}

	projPath, err := utils.GetLocalPath(projPath)
	if err == nil {
		if !utils.FileExists(projPath) {
			os.MkdirAll(projPath, 0755)
		}
	}

	return projPath, err

}

func isSupported(val string, supported []string) bool {
	for i := range supported {
		if supported[i] == val {
			return true
		}
	}
	return false
}

func setLanguage(ctx *cli.Context, devpacks *packs.DevPacks, dir string) (*packs.DevPack, error) {
	supported, err := devpacks.List()
	if err != nil {
		return nil, err
	}

	// Do not prompt if input is valid
	lang := ctx.String("lang")
	if isSupported(lang, supported) {
		return devpacks.Load(lang)
	}

	// Set guestimate as default
	lang = guessLang(dir)

	prompt := "Language"
	lang = promptForSupported(prompt, supported, lang)

	devpack, err := devpacks.Load(lang)
	if err == nil {
		err = ctx.Set("lang", devpack.Name+":"+devpack.DefaultVersion)
	}

	return devpack, err
}

func setRepoOwner(ctx *cli.Context, defRepoOwner, dir string) error {
	var err error

	var repoOwner string
	utils.PromptUntilNoError("Repo owner ["+defRepoOwner+"]: ", os.Stdout, os.Stdin, func(db []byte) error {
		repoOwner = string(db)
		if repoOwner == "" {
			repoOwner = defRepoOwner
		}
		return nil
	})

	err = ctx.Set(vars.VcsRepoOwner, repoOwner)
	return err
}

func guessLang(dir string) string {
	fts := analysis.BuildFileTypeSpread(dir)
	highest := fts.Highest()
	if highest != nil && highest.Percent > 50 {
		return strings.ToLower(highest.Language)
	}
	return ""
}

func promptComps(name, lang string, pks *packs.Packs) (*thrapb.Stack, error) {
	c := &thrap.BasicStackConfig{
		Name:     name,
		Language: thrapb.LanguageID(lang),
	}

	var err error
	// ds, ws := prompt(conf)
	c.WebServer, err = promptPack(pks.Web(), "Web Server")
	if err != nil {
		return nil, err
	}
	c.DataStore, err = promptPack(pks.Datastore(), "Data Store")
	if err != nil {
		return nil, err
	}

	return thrap.NewBasicStack(c, pks)
}

func promptPack(wp *packs.BasePacks, prompt string) (string, error) {
	list, err := wp.List()
	if err != nil {
		return "", err
	}
	supported := append(list, "none")
	return promptForSupported(prompt, supported, "none"), nil
}

//
// func prompt(conf *config.Config) (string, string) {
// 	var (
// 		prompt    string
// 		supported []string
// 	)
//
// 	// Prompt for datastore
// 	supported = make([]string, 0, len(conf.DataStores)+1)
// 	for k := range conf.DataStores {
// 		supported = append(supported, k)
// 	}
// 	prompt = "Data store"
// 	ds := promptForSupported(prompt, supported, "none")
//
// 	// Prompt for webserver
// 	supported = make([]string, 0, len(conf.WebServers)+1)
// 	for k := range conf.WebServers {
// 		supported = append(supported, k)
// 	}
// 	prompt = "Web server"
// 	ws := promptForSupported(prompt, supported, "none")
//
// 	return ds, ws
// }

func setupGit(projName, repoOwner, projPath, remoteAddr string) (vcs.VCS, *git.Repository, error) {
	vcsp := vcs.NewGitVCS()

	rr := &vcs.Repository{Name: projName}
	opt := vcs.Option{
		Path:   projPath,
		Remote: vcs.DefaultGitRemoteURL(remoteAddr, repoOwner, projName),
	}

	resp, err := vcsp.Create(rr, opt)
	if err != nil {
		return vcsp, nil, err
	}

	return vcsp, resp.(*git.Repository), nil

}
