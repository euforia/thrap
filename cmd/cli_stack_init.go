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
	"github.com/euforia/thrap/config"
	"github.com/euforia/thrap/consts"
	"github.com/euforia/thrap/devpack"
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
		Usage:     "Initialize new project",
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

			conf, err := config.ParseFile("./etc/config.hcl")
			if err != nil {
				return err
			}

			projPath, err := setupProjPath(ctx)
			if err != nil {
				return err
			}

			// Project name
			projName := ctx.String("name")
			if len(projName) == 0 {
				projName = filepath.Base(projPath)
			}

			devpacks := devpack.NewDevPacks("./etc/packs/dev")
			// Set language from input or otherwise and other related params
			_, err = setLanguage(ctx, devpacks, projPath)
			if err != nil {
				return err
			}

			thrapConf, err := thrap.ReadGlobalConfig()
			if err != nil {
				return err
			}

			if err = setRepoOwner(ctx, thrapConf.VCS.Username, projPath); err != nil {
				return err
			}

			//projLang := ctx.String("lang")
			repoOwner := ctx.String(vars.VcsRepoOwner)

			// Local project setup
			err = thrap.ConfigureProjectDir(projName, repoOwner, projPath)
			if err != nil {
				return err
			}

			mfile := filepath.Join(projPath, consts.DefaultManifestFile)
			if utils.FileExists(mfile) {
				// TODO: ???
				return fmt.Errorf("manifest %s already exists", consts.DefaultManifestFile)
			}

			scopeVars, _ := thrap.LoadGlobalScopeVars()
			// Need the project dir configured before scope vars can be loaded
			projVars, _ := thrap.LoadProjectScopeVars(projPath)
			scopeVars = vars.MergeScopeVars(scopeVars, projVars)

			gitRemoteAddr := scopeVars[vars.VcsAddr].Value.(string)
			vcsp, gitRepo, err := setupGit(projName, repoOwner, projPath, gitRemoteAddr)
			if err != nil {
				return err
			}

			// Prompt for missing
			stack := promptComps(projName, ctx.String("lang"), conf)
			errs := stack.Validate()
			if errs != nil {
				return utils.FlattenErrors(errs)
			}

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

func isLangSupported(val string, supported []string) bool {
	for i := range supported {
		if supported[i] == val {
			return true
		}
	}
	return false
}

func setLanguage(ctx *cli.Context, devpacks *devpack.DevPacks, dir string) (*devpack.DevPack, error) {
	supported, err := devpacks.List()
	if err != nil {
		return nil, err
	}

	// Do not prompt if input is valid
	lang := ctx.String("lang")
	if isLangSupported(lang, supported) {
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
	// if ctx.String("lang") == "go" {
	// 	defRepoOwner = filepath.Base(filepath.Dir(dir))
	// }

	var repoOwner string
	promptUntilNoError("Repo owner ["+defRepoOwner+"]: ", os.Stdout, os.Stdin, func(db []byte) error {
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

func promptComps(name, lang string, conf *config.Config) *thrapb.Stack {
	ds, ws := prompt(conf)
	st := thrap.NewBasicStack(thrapb.LanguageID(lang), name, ds, ws, conf)
	fmt.Println()
	return st
}

func prompt(conf *config.Config) (string, string) {
	var (
		prompt    string
		supported []string
	)

	// Prompt for datastore
	supported = make([]string, 0, len(conf.DataStores)+1)
	for k := range conf.DataStores {
		supported = append(supported, k)
	}
	prompt = "Data store"
	ds := promptForSupported(prompt, supported, "none")

	// Prompt for webserver
	supported = make([]string, 0, len(conf.WebServers)+1)
	for k := range conf.WebServers {
		supported = append(supported, k)
	}
	prompt = "Web server"
	ws := promptForSupported(prompt, supported, "none")

	return ds, ws
}

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
