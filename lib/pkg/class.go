package pkg

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

type PackageV1 struct {
	PackName string
	Priority uint64
}

type PackagesV1 []PackageV1

func (p PackagesV1) Len() int {
	return len(p)
}

func (p PackagesV1) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p PackagesV1) Less(i, j int) bool {
	return p[i].Priority < p[j].Priority
}
