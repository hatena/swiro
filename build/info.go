package build

type Info struct {
	Name    string
	Version string
}

var (
	version string
	name    string
)

func GetInfo() *Info {
	return &Info{
		Name:    name,
		Version: version,
	}
}
