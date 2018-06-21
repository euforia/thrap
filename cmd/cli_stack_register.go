package main

import (
	"context"
	"fmt"
	"os"

	"github.com/euforia/hclencoder"
	"github.com/euforia/thrap"
	"github.com/euforia/thrap/thrapb"
	"github.com/euforia/thrap/utils"
	"google.golang.org/grpc"
	"gopkg.in/urfave/cli.v2"
)

func commandStackRegister() *cli.Command {
	return &cli.Command{
		Name:  "register",
		Usage: "Register project",
		Action: func(ctx *cli.Context) error {

			stack, err := thrap.LoadManifest("")
			if err != nil {
				return err
			}
			if errs := stack.Validate(); errs != nil {
				return utils.FlattenErrors(errs)
			}

			taddr := ctx.String("thrap-addr")
			if taddr == "" {
				return errRemoteRequired
			}
			cc, err := grpc.Dial(taddr, grpc.WithInsecure())
			if err != nil {
				return err
			}

			tclient := thrapb.NewThrapClient(cc)
			resp, err := tclient.RegisterStack(context.Background(), stack)
			if err != nil {
				return err
			}

			m := map[string]*thrapb.Stack{
				"stack": resp,
			}

			b, _ := hclencoder.Encode(m)
			fmt.Printf("\n%s\n", b)
			return nil
		},
	}
}

func commandStackValidate() *cli.Command {
	return &cli.Command{
		Name:      "validate",
		Usage:     "Validate a manifest",
		ArgsUsage: "[path to manifest]",
		Action: func(ctx *cli.Context) error {
			mfile := ctx.Args().Get(0)
			mf, err := thrap.LoadManifest(mfile)
			if err != nil {
				return err
			}

			errs := mf.Validate()
			if errs == nil {
				writeHCLManifest(mf, os.Stdout)
			} else {
				err = utils.FlattenErrors(errs)
			}

			return err
		},
	}
}

// func loadRegistryProvider(varmap map[string]string) (registry.Registry, error) {
// 	conf := registry.DefaultConfig()
// 	conf.Provider = varmap[vars.RegistryID]
// 	conf.Conf["addr"] = varmap[vars.RegistryAddr]
//
// 	return registry.New(conf)
// }

// func loadVcsProvider(allvars map[string]string) (vcs.VCS, error) {
// 	creds := make(map[string]string, 1)
// 	for k, v := range allvars {
// 		if strings.HasPrefix(k, vars.VcsCreds) {
// 			key := strings.TrimPrefix(k, vars.VcsCreds+".")
// 			creds[key] = v
// 		}
// 	}
//
// 	return vcs.New(&vcs.Config{Provider: allvars[vars.VcsID], Conf: creds})
// }

// func loadVars(path string) (map[string]string, error) {
// 	projpath, err := thrap.GetLocalPath(path)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	all, err := thrap.LoadGlobalVars()
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	pvars, err := thrap.LoadVars(projpath)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return thrap.MergeVars(all, pvars), nil
// }
