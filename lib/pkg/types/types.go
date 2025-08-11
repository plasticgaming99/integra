package types

// minimum struct
type Pkg struct {
	PkgName string
	Version string
	Release int
}

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

func PackInfoToPkg(pkginfo Packinfo) Pkg {
	return Pkg{
		PkgName: pkginfo.Packagename,
		Version: pkginfo.Version,
		Release: pkginfo.Release,
	}
}
