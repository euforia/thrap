package core

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"text/tabwriter"
	"time"

	"github.com/euforia/pseudo/scope"

	"github.com/euforia/thrap/pkg/metrics"
	"github.com/euforia/thrap/pkg/provider"
	"github.com/euforia/thrap/thrapb"
)

func printScopeVars(scopeVars scope.Variables) {
	fmt.Printf("\nScope:\n\n")
	for _, name := range scopeVars.Names() {
		fmt.Println(" ", name)
	}
	fmt.Println()
}

// getBuildImageTags returns tags that should be applied to a given image build. If a
// registry config is provided, names are generated accordingly
func getBuildImageTags(sid string, comp *thrapb.Component, rconf *provider.Config) []string {
	base := filepath.Join(sid, comp.ID)

	out := []string{}
	if rconf != nil && len(rconf.Addr) > 0 {
		// remote
		rbase := filepath.Join(rconf.Addr, base)
		out = append(out, rbase)
		if len(comp.Version) > 0 {
			out = append(out, rbase+":"+comp.Version)
		}
	} else {
		// local
		out = []string{base}
		if len(comp.Version) > 0 {
			out = append(out, base+":"+comp.Version)
		}
	}
	return out
}

func mapHasErrors(m map[string]error) bool {
	for _, v := range m {
		if v != nil {
			return true
		}
	}
	return false
}

func printPublishResults(results map[string]error) {
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.StripEscape)
	fmt.Fprintf(tw, " \tArtifact\tStatus\tDetails\n")
	fmt.Fprintf(tw, " \t--------\t------\t-------\n")
	for image, err := range results {
		if err != nil {
			fmt.Fprintf(tw, " \t%s\tfailed\t%v\n", image, err)
		} else {
			fmt.Fprintf(tw, " \t%s\tsucceeded\t\n", image)
		}

	}
	tw.Flush()
	fmt.Println()
}

func printBuildResults(stack *thrapb.Stack, results map[string]*CompBuildResult, w io.Writer) {
	w.Write([]byte("\n"))
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', tabwriter.StripEscape)
	fmt.Fprintf(tw, " \tComponent\tArtifact\tStatus\tDetails\n")
	fmt.Fprintf(tw, " \t---------\t--------\t------\t-------\n")
	for k, r := range results {
		var (
			status string
			msg    string
			art    string
		)

		if r.Error == nil {
			status = "succeeded"
			art = stack.ArtifactName(k) + ":" + stack.Components[k].Version
		} else {
			status = "failed"
			msg = r.Error.Error()
		}

		fmt.Fprintf(tw, " \t%s\t%s\t%s\t%s\n", k, art, status, msg)
	}
	tw.Flush()
	w.Write([]byte("\n"))
}

func printBuildStats(bld *stackBuilder, total, pub *metrics.Runtime) {
	s := bld.ServiceTime()
	b := bld.BuildTime()
	results := bld.Results()

	fmt.Printf("\n  Timing:\n\n   Service:\t%v\n", s.Duration(time.Millisecond))
	fmt.Printf("   Build:\t%v\n", b.Duration(time.Millisecond))
	for k, v := range results {
		fmt.Printf("     %s:\t%v\n", k, v.Runtime.Duration(time.Millisecond))
	}
	fmt.Printf("   Publish:\t%v\n\n", pub.Duration(time.Millisecond))
	fmt.Printf("   Total:\t%v\n", total.Duration(time.Millisecond))
}

func printScopeVarsWithVals(svars scope.Variables) {

	s := make([]string, 0, len(svars))
	for k := range svars {
		s = append(s, k)
	}
	sort.Strings(s)

	fmt.Printf("\nScope:\n\n")
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.StripEscape)
	for _, k := range s {
		v := svars[k]
		fmt.Fprintf(tw, " \t%s\t%v\n", k, v.Value)
	}
	tw.Flush()
}
