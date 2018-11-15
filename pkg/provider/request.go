package provider

import "github.com/euforia/thrap/pkg/pb"

// Request is a raw request for a provider
type Request struct {
	Project    pb.Project
	Deployment pb.Deployment
	// Unmarshalled Deployment.Spec
	// Spec interface{}
}
