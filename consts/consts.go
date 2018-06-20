package consts

// WorkDir is the global working directory name
const (
	// WorkDir is the name of the thrap working directory global and local
	WorkDir = ".thrap"
	// ConfigFile is the default config filename
	ConfigFile = "config.hcl"
	// CredsFile is the default credentials filename
	CredsFile = "creds.hcl"
	// KeyFile is the keypair file
	KeyFile = "ecdsa256"
	// EnvVarVersion is the env. var. name for the version injected by thrap
	EnvVarVersion = "APP_VERSION"
)

const (
	// DefaultManifestFile is the default manifest filename
	DefaultManifestFile = "thrap.hcl"
	// DefaultWorkDir is the default working directory set if not provided
	DefaultWorkDir = "/"
)

const (
	// DefaultReadmeFile is the readme filename
	DefaultReadmeFile = "README.md"
	// DefaultSecretsFile is the secrets filename
	DefaultSecretsFile = "secrets.hcl"
	DefaultMakefile    = "Makefile"
	DefaultEnvFile     = ".env"
	DefaultDockerFile  = "dockerfile"
	// DefaultBuildContext is the container build context
	DefaultBuildContext = "."
	DefaultWebCompID    = "www"
	DefaultAPICompID    = "api"
	DefaultDSCompID     = "db"
)
