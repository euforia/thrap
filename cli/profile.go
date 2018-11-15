package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/euforia/hclencoder"
	"github.com/euforia/thrap/pkg/pb"
	"github.com/euforia/thrap/pkg/storage"

	"github.com/euforia/thrap/utils"
	"gopkg.in/urfave/cli.v2"
)

func commandProfile() *cli.Command {
	return &cli.Command{
		Name:      "profile",
		Usage:     "Iteract with profiles",
		ArgsUsage: "<profile>",
		Action: func(ctx *cli.Context) error {
			ppath, err := utils.GetLocalPath("")
			if err != nil {
				return err
			}

			profs, err := storage.LoadHCLFileProfileStorage(ppath)
			if err != nil {
				return err
			}

			var (
				profIn  = ctx.Args().Get(0)
				display interface{}
			)

			if profIn == "" {
				display = profs.List()
				os.Stdout.Write([]byte("\n"))
			} else {
				kv := strings.Split(profIn, "=")
				if len(kv) < 1 {
					cli.ShowCommandHelpAndExit(ctx, "profile", 1)
				}
				id := kv[0]

				prof, _ := profs.Get(id)
				if prof == nil {
					return storage.ErrProfileNotFound
				}
				display = map[string]*pb.Profile{id: prof}
				fmt.Println()
			}

			b, _ := hclencoder.Encode(&display)
			fmt.Printf("%s\n", b)

			return nil
		},
	}
}
