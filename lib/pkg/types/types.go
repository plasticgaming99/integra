package types

// much basic packinfo
type Packinfo struct {
	Packagename  string
	Version      string
	Release      int
	License      string
	Architecture string
	Description  string
	Url          string
	Depends      []string
	Optdeps      []string
	Builddeps    []string
	Conflicts    []string
	Provides     []string
}
