package environment

const (
	dev = "0.0.0-dev"
)

var (
	version      = dev             // Version of this binary
	gitVersion   = "v0.0.0-master" // nolint:unused
	gitCommit    = "unknown"       // nolint:unused // sha1 from git, output of $(git rev-parse HEAD)
	gitTreeState = "unknown"       // nolint:unused // state of git tree, either "clean" or "dirty"
	buildDate    = "unknown"       // nolint:unused // build date in ISO8601 format, output of $(date -u +'%Y-%m-%dT%H:%M:%SZ')
)

func Version() string {
	return version
}

// IsDev returns true if the version is dev.
func IsDev() bool {
	return version == dev
}

// GitVersion returns the git version of this binary.
func GitVersion() string {
	return gitVersion
}

// GitCommit returns the git commit of this binary.
func GitCommit() string {
	return gitCommit
}

// GitTreeState returns the git tree state of this binary.
func GitTreeState() string {
	return gitTreeState
}

// BuildDate returns the build date of this binary.
func BuildDate() string {
	return buildDate
}
