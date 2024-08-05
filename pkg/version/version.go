package version

var (
	version string
	commit  string
	date    string
)

type Info struct {
	Version string
	Date    string
	Commit  string
}

func GetVersionInfo() Info {
	return Info{
		Version: version,
		Date:    date,
		Commit:  commit,
	}
}

func GetVersion() string {
	return version
}
