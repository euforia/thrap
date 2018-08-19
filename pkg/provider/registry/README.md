# registry

This package implements a unified interface to interact with container registries.

Current supported registries are:

- ECR
- Docker Hub
- Local Docker runtime

The registry interface is as follows:

```golang
type Registry interface {
    ID() string
    // Initialize the registry provider
    Init(conf *Config) error
    // Create a new repository
    CreateRepo(string) (interface{}, error)
    // Get repo info
    GetRepo(string) (interface{}, error)
    // Get image manifest
    GetImageManifest(name, tag string) (interface{}, error)
    // Name of the image with the registry. Needed for deployments
    ImageName(string) string
    // Returns a docker AuthConfig
    GetAuthConfig() (types.AuthConfig, error)
}
```
