package thrapb

import (
	"fmt"
	"io"
	"text/tabwriter"
)

type ActionResult struct {
	Action   string
	Resource string
	Data     interface{}
	Error    error
}

type ActionsResults map[string][]*ActionResult

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

// // Action represents any noteworthy command, transaction etc.
// type Action struct {
// 	// Name of the action
// 	Name string
// 	// Type of resource
// 	Resource string
// 	// Resource identifier
// 	Identifier string
// }

// // NewAction returns an action with the given parameters
// func NewAction(name, rsrc, id string) *Action {
// 	return &Action{
// 		Name:       name,
// 		Resource:   rsrc,
// 		Identifier: id,
// 	}
// }

// func (a *Action) String() string {
// 	return a.Resource + " " + a.Identifier + " " + a.Name
// }

// // ActionReport holds an execution report for a given action
// type ActionReport struct {
// 	Action *Action
// 	Data   interface{}
// 	Error  error
// }

// // HasError returns true if the ActionReport contains an error
// func (ar *ActionReport) HasError() bool {
// 	return ar.Error != nil
// }
