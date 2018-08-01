package core

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/euforia/pseudo/scope"
	"github.com/euforia/thrap/config"
	"github.com/euforia/thrap/metrics"
	"github.com/euforia/thrap/thrapb"
)

func printScopeVars(scopeVars scope.Variables) {
	fmt.Printf("\nScope:\n\n")
	for _, name := range scopeVars.Names() {
		fmt.Println(" ", name)
	}
	fmt.Println()
}

// getBuildImageTags returns tags that should be applied to a given image build
func getBuildImageTags(sid string, comp *thrapb.Component, rconf *config.RegistryConfig) []string {
	base := filepath.Join(sid, comp.ID)
	out := []string{base}
	if len(comp.Version) > 0 {
		out = append(out, base+":"+comp.Version)
	}

	// rconf := bldr.conf //.GetDefaultRegistry()
	if rconf != nil && len(rconf.Addr) > 0 {
		rbase := filepath.Join(rconf.Addr, base)
		out = append(out, rbase)
		if len(comp.Version) > 0 {
			out = append(out, rbase+":"+comp.Version)
		}
	}
	return out
}
func printBuildResults(stack *thrapb.Stack, results map[string]*CompBuildResult, w io.Writer) {
	w.Write([]byte("\n"))
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', tabwriter.StripEscape)
	fmt.Fprintf(tw, "Component\tArtifact\tStatus\tDetails\n")
	fmt.Fprintf(tw, "---------\t--------\t------\t-------\n")
	for k, r := range results {
		var (
			status string
			msg    string
			art    string
		)
		if r.Error == nil {
			comp := stack.Components[k]
			status = "success"
			art = comp.Name + ":" + comp.Version
		} else {
			status = "fail"
			msg = r.Error.Error()
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", k, art, status, msg)
	}
	tw.Flush()
	w.Write([]byte("\n"))
}

func printArtifacts(stack *thrapb.Stack, rconf *config.RegistryConfig) {
	fmt.Printf("\nArtifacts:\n\n")
	for _, comp := range stack.Components {
		if !comp.IsBuildable() {
			continue
		}

		names := getBuildImageTags(stack.ID, comp, rconf)
		fmt.Println(strings.Join(names, "\n"))
		fmt.Println()
	}
}

func printBuildStats(bld *stackBuilder, total, pub *metrics.Runtime) {
	s := bld.ServiceTime()
	b := bld.BuildTime()
	results := bld.Results()

	fmt.Printf("Timing:\n\n  Service:\t%v\n", s.Duration(time.Millisecond))
	fmt.Printf("  Build:\t%v\n", b.Duration(time.Millisecond))
	for k, v := range results {
		fmt.Printf("    %s:\t%v\n", k, v.Runtime.Duration(time.Millisecond))
	}
	fmt.Printf("  Publish:\t%v\n\n", pub.Duration(time.Millisecond))
	fmt.Printf("  Total:\t%v\n", total.Duration(time.Millisecond))
}
