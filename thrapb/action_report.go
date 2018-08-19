package thrapb

import (
	"fmt"
	"io"
	"text/tabwriter"
)

// ActionResult holds the result of a given action
type ActionResult struct {
	Action   string
	Resource string
	Data     interface{}
	Error    error
}

// ActionsResults is a collection of printable results
type ActionsResults map[string][]*ActionResult

// Print prints the formatted results
func (results ActionsResults) Print(w io.Writer) {

	for k, slice := range results {
		fmt.Printf("%s\n\n", k)
		tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', tabwriter.StripEscape)
		fmt.Fprintf(tw, " \tID\tStatus\tDetails\n")
		fmt.Fprintf(tw, " \t--\t------\t-------\n")
		for _, result := range slice {
			if result.Error == nil {
				fmt.Fprintf(tw, " \t%s\tok\t%v\n", result.Resource, result.Data)
			} else {
				fmt.Fprintf(tw, " \t%s\terror\t%v\n", result.Resource, result.Error)
			}
		}
		tw.Flush()
		fmt.Println()
	}

}
