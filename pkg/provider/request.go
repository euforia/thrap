package provider

import "github.com/euforia/thrap/thrapb"

// Request is a raw request for a provider
type Request struct {
	Project    thrapb.Project
	Deployment thrapb.Deployment
	// Unmarshalled Deployment.Spec
	// Spec interface{}
}
