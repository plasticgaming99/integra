package types

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

func ReadPackinfo(in io.Reader) (pkinfo Packinfo) {
	scan := bufio.NewScanner(in)
	// reuse for cut
	var (
		key  string
		val  string
		cont bool

		err error
	)
	for scan.Scan() {
		key, val, cont = strings.Cut(scan.Text(), " = ")
		if !cont {
			continue
		}
		switch key {
		case "package":
			pkinfo.Packagename = val
		case "version":
			pkinfo.Version = val
		case "release":
			pkinfo.Release, err = strconv.Atoi(val)
			// fallback, not critical
			if err != nil {
				pkinfo.Release = 0
			}
		case "license":
			pkinfo.License = val
		case "architecture":
			pkinfo.Architecture = val
		case "description":
			pkinfo.Description = val
		case "url":
			pkinfo.Url = val
		case "depends":
			pkinfo.Depends = append(pkinfo.Depends, val)
		case "optdepends":
			pkinfo.Optdeps = append(pkinfo.Optdeps, val)
		case "builddeps":
			pkinfo.Builddeps = append(pkinfo.Builddeps, val)
		case "conflicts":
			pkinfo.Conflicts = append(pkinfo.Conflicts, val)
		case "provides":
			pkinfo.Provides = append(pkinfo.Provides, val)
		default:
			continue
		}

	}
	return
}
