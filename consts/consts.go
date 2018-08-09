package consts

// WorkDir is the global working directory name
const (
	// WorkDir is the name of the thrap working directory global and local
	WorkDir = ".thrap"
	// ConfigFile is the default config filename
	ConfigFile = "config.hcl"
	// CredsFile is the default credentials filename
	CredsFile = "creds.hcl"
	// IdentityFile is the identity filename
	IdentityFile = "identity.hcl"
	// ProfilesFile is the profiles filename
	ProfilesFile = "profiles.hcl"
	// KeyFile is the keypair file
	KeyFile = "ecdsa256"
	// EnvVarVersion is the env. var. name for the version injected by thrap
	EnvVarVersion = "STACK_VERSION"
	// PacksDir is the directory name where packs are stored
	PacksDir = "packs"
)

const (
	// DefaultManifestFile is the default manifest filename
	DefaultManifestFile = "thrap.yml"
	// DefaultWorkDir is the default container working directory set if not
	// provided
	DefaultWorkDir = "/"
	// DefaultBuildContext is the container build context
	DefaultBuildContext = "."
	//
	DefaultDataDir = "~/.thrap"
)

const (
	// DefaultReadmeFile is the readme filename
	DefaultReadmeFile = "README.md"
	// DefaultSecretsFile is the secrets filename
	DefaultSecretsFile = "secrets.hcl"
	// DefaultSecretsFileFormat is the secrets format
	DefaultSecretsFileFormat = "hcl"
	// DefaultMakefile is the makefile filename
	DefaultMakefile = "Makefile"
	// DefaultEnvFile is the env filename
	DefaultEnvFile = ".env"
	// DefaultDockerFile is the docker filename
	DefaultDockerFile = "dockerfile"
)

const (
	// CompVarPrefixKey is the prefix for all component variables
	CompVarPrefixKey = "comp"
	// DepVarPrefixKey is the prefix for all dep variables
	DepVarPrefixKey = "dep"
)

const (
	// DefaultWebCompID is the default web Component id
	DefaultWebCompID = "www"
	// DefaultAPICompID is the default api Component id
	DefaultAPICompID = "api"
	// DefaultDSCompID is the default datastore Component id
	DefaultDSCompID = "db"
)
