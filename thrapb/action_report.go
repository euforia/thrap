package thrapb

// Action represents any noteworthy command, transaction etc.
type Action struct {
	// Name of the action
	Name string
	// Type of resource
	Resource string
	// Resource identifier
	Identifier string
}

// NewAction returns an action with the given parameters
func NewAction(name, rsrc, id string) *Action {
	return &Action{
		Name:       name,
		Resource:   rsrc,
		Identifier: id,
	}
}

func (a *Action) String() string {
	return a.Resource + " " + a.Identifier + " " + a.Name
}

// ActionReport holds an execution report for a given action
type ActionReport struct {
	Action *Action
	Data   interface{}
	Error  error
}

// HasError returns true if the ActionReport contains an error
func (ar *ActionReport) HasError() bool {
	return ar.Error != nil
}
