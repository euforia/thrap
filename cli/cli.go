package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"google.golang.org/grpc"

	"github.com/euforia/thrap/consts"
	"github.com/euforia/thrap/store"
	"github.com/euforia/thrap/utils"
	"github.com/euforia/thrap/vars"

	"github.com/pkg/errors"

	"github.com/euforia/thrap/config"
	"github.com/euforia/thrap/core"
	"github.com/euforia/thrap/thrapb"
	"gopkg.in/urfave/cli.v2"
)

var (
	errThrapAddrRequired = errors.New("--thrap-addr required")
	errNotConfigured     = errors.New("thrap not configured. Try running 'thrap configure'")
)

// NewCLI returns a new command line app
func NewCLI(version string) *cli.App {
	cli.VersionPrinter = func(ctx *cli.Context) {
		fmt.Println(ctx.App.Version)
	}

	app := &cli.App{
		Name:     "thrap",
		HelpName: "thrap",
		Version:  version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "thrap-addr",
				Usage:   "thrap registry address",
				EnvVars: []string{"THRAP_ADDR"},
			},
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "Debug mode",
			},
		},
		Commands: []*cli.Command{
			commandConfigure(),
			commandIdentity(),
			commandAgent(),
			commandStack(),
			commandPack(),
			commandVersion(),
		},
	}

	app.HideVersion = true

	return app
}

func commandVersion() *cli.Command {
	return &cli.Command{
		Name:  "version",
		Usage: "Show version",
		Action: func(ctx *cli.Context) error {
			fmt.Println(ctx.App.Version)
			return nil
		},
	}
}

func commandConfigure() *cli.Command {
	return &cli.Command{
		Name:  "configure",
		Usage: "Configure global settings",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:   vars.VcsID,
				Usage:  "version control `provider`",
				Value:  "github",
				Hidden: true,
			},
			&cli.StringFlag{
				Name:  vars.VcsUsername,
				Usage: "version control `username`",
			},
			&cli.StringFlag{
				Name:  "data-dir",
				Usage: "data `directory`",
				Value: "~/" + consts.WorkDir,
			},
			&cli.BoolFlag{
				Name:  "no-prompt",
				Usage: "do not prompt for input",
			},
		},
		Action: func(ctx *cli.Context) error {
			opts := core.ConfigureOptions{
				VCS: &config.VCSConfig{
					ID:       ctx.String(vars.VcsID),
					Username: ctx.String(vars.VcsUsername),
				},
				DataDir:  ctx.String("data-dir"),
				NoPrompt: ctx.Bool("no-prompt"),
			}

			// Only configures things that are not configured
			return core.ConfigureGlobal(opts)
		},
	}
}

func newThrapClient(ctx *cli.Context) (thrapb.ThrapClient, error) {
	// Check remote addr
	remoteAddr := ctx.String("thrap-addr")
	if remoteAddr == "" {
		return nil, errThrapAddrRequired
	}

	cc, err := grpc.Dial(remoteAddr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return thrapb.NewThrapClient(cc), nil
}

func writeJSON(v interface{}) {
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Printf("%s\n", b)
}

func loadCore(ctx *cli.Context) (*core.Core, error) {
	lpath, err := utils.GetLocalPath("")
	if err != nil {
		return nil, err
	}

	conf := &core.Config{DataDir: consts.DefaultDataDir}
	if ctx.Bool("debug") {
		conf.Logger = core.DefaultLogger(os.Stdout)
	}

	// Load project configs
	conf.ThrapConfig, err = config.ReadProjectConfig(lpath)
	if err != nil {
		return nil, err
	}

	conf.Creds, err = config.ReadProjectCredsConfig(lpath)
	if err != nil {
		return nil, err
	}

	return core.NewCore(conf)
}

func loadProfile(ctx *cli.Context) (*store.HCLFileProfileStorage, *thrapb.Profile, error) {
	lpath, err := utils.GetLocalPath("")
	if err != nil {
		return nil, nil, err
	}

	// Load profiles
	profs, err := store.LoadHCLFileProfileStorage(lpath)
	if err != nil {
		return nil, nil, err
	}

	// Load request profile
	profName := ctx.String("profile")
	prof := profs.Get(profName)
	if prof == nil {
		return profs, nil, fmt.Errorf("profile not found: %s", profName)
	}

	fmt.Printf("Profile: %s\n\n", profName)
	return profs, prof, nil
}
