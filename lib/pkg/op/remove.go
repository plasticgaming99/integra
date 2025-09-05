package op

import (
	"fmt"
	"os"

	"github.com/plasticgaming99/integra/lib/db/localdb"
	"github.com/plasticgaming99/integra/lib/integrity"
	"github.com/plasticgaming99/integra/lib/pkg/types"
)

// remove function, it does not care about dependency too
func Remove(pkg types.Pkg, rootdir string, ldb localdb.LocalDB) error {
	rd, err := ldb.GetFile(pkg, ".INTEGRITY")
	if err != nil {
		return fmt.Errorf("error getting INTEGRITY for package %s", localdb.PkgToDirname(pkg))
	}
	intg := integrity.Parse(rd)
	for _, intg := range intg.Files {
		os.Remove(intg.Filepath)
	}
	ldb.UnregisterPkg(pkg)
	return nil
}
