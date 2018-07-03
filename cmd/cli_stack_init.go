package main

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/urfave/cli.v2"

	"github.com/euforia/thrap/analysis"
	"github.com/euforia/thrap/asm"
	"github.com/euforia/thrap/config"
	"github.com/euforia/thrap/consts"
	"github.com/euforia/thrap/core"
	"github.com/euforia/thrap/manifest"
	"github.com/euforia/thrap/packs"
	"github.com/euforia/thrap/thrapb"
	"github.com/euforia/thrap/utils"
	"github.com/euforia/thrap/vars"
	"github.com/pkg/errors"
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
				Name:   vars.VcsID,
				Usage:  "version control `provider` (experimental)",
				Value:  "github",
				Hidden: true,
			},
			&cli.StringFlag{
				Name:  vars.VcsRepoOwner,
				Usage: "source code repo `owner`",
			},
		},
		Action: func(ctx *cli.Context) error {
			coreConf := &core.Config{
				PacksDir: defaultPacksDir,
			}

			cr, err := core.NewCore(coreConf)
			if err != nil {
				return err
			}

			// Create project directory if not found
			projPath, err := setupProjPath(ctx)
			if err != nil {
				return err
			}

			// Project name
			projName := ctx.String("name")
			if len(projName) == 0 {
				projName = filepath.Base(projPath)
			}

			mfile := filepath.Join(projPath, consts.DefaultManifestFile)
			if utils.FileExists(mfile) {
				// TODO: ???
				return fmt.Errorf("manifest %s already exists", consts.DefaultManifestFile)
			}

			pks := cr.Packs()
			// Set language from input or otherwise and other related params
			_, err = setLanguage(ctx, pks.Dev(), projPath)
			if err != nil {
				return err
			}

			gconf := cr.Config()
			vcsID := ctx.String(vars.VcsID)
			defaultVCS := gconf.VCS[vcsID]
			repoOwner := setRepoOwner(ctx, defaultVCS.ID, defaultVCS.Username)

			// Local project setup
			opts := core.ConfigureOptions{
				DataDir: projPath,
				VCS: &config.VCSConfig{
					ID: defaultVCS.ID,
					Repo: &config.VCSRepoConfig{
						Name:  projName,
						Owner: repoOwner,
					},
				},
			}

			// Prompt for missing
			bsc, err := promptComps(projName, ctx.String("lang"), pks)
			if err != nil {
				return err
			}

			fmt.Println()

			stm := cr.Stack()
			stack, err := stm.Init(bsc, opts)
			if err != nil {
				return err
			}

			return manifest.WriteManifest(stack, os.Stdout)
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
	if !isSupported(lang, supported) {

		// Set guestimate as default
		lang = analysis.EstimateLanguage(dir)

		prompt := "Language"
		lang = promptForSupported(prompt, supported, lang)
	}

	devpack, err := devpacks.Load(lang)
	if err == nil {
		err = ctx.Set("lang", devpack.Name+":"+devpack.DefaultVersion)
	}

	return devpack, err
}

func setRepoOwner(ctx *cli.Context, vcsID, defRepoOwner string) string {
	//var err error

	var repoOwner string
	prompt := vcsID + " repo owner [" + defRepoOwner + "]: "
	utils.PromptUntilNoError(prompt, os.Stdout, os.Stdin, func(db []byte) error {
		repoOwner = string(db)
		if repoOwner == "" {
			if defRepoOwner == "" {
				return errors.New("repo owner required")
			}
			repoOwner = defRepoOwner
		}
		return nil
	})

	return repoOwner
}

func promptComps(name, lang string, pks *packs.Packs) (*asm.BasicStackConfig, error) {
	c := &asm.BasicStackConfig{
		Name:     name,
		Language: thrapb.LanguageID(lang),
	}

	var err error
	c.WebServer, err = promptPack(pks.Web(), "Web Server")
	if err != nil {
		return nil, err
	}
	c.DataStore, err = promptPack(pks.Datastore(), "Data Store")
	if err != nil {
		return nil, err
	}

	return c, nil
}

func promptPack(wp *packs.BasePacks, prompt string) (string, error) {
	list, err := wp.List()
	if err != nil {
		return "", err
	}
	supported := append(list, "none")
	return promptForSupported(prompt, supported, "none"), nil
}
