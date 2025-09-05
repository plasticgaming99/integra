package op

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/plasticgaming99/integra/lib/db/localdb"
	"github.com/plasticgaming99/integra/lib/integrity"
	"github.com/plasticgaming99/integra/lib/pkg/types"
	"github.com/zeebo/blake3"
)

func Check(pkg types.Pkg, rootdir string, ldb localdb.LocalDB) (bool, error) {
	rd, err := ldb.GetFile(pkg, localdb.IntegrityFile)
	if err != nil {
		return false, fmt.Errorf("error getting INTEGRITY file from localdb: %w", err)
	}
	defer rd.Close()
	bufrd := bufio.NewReader(rd)
	intg := integrity.Parse(bufrd)
	for _, in := range intg.Files {
		if in.Filepath == "/.PACKAGE" {
			continue
		}
		path := filepath.Join(rootdir, in.Filepath)
		stat, err := os.Lstat(path)
		if err != nil {
			fmt.Println("error opening file: ", err)
			continue
		}

		if stat.IsDir() {
			if stat.Mode() == in.FileMode.Perm() {
				continue
			}
		} else {
			file, errf := os.Open(path)
			sum := []byte{}
			linksto := ""
			if errf != nil {
				p, er := os.Readlink(path)
				if er != nil {
					file.Close()
					fmt.Println("error reading link:", er)
					fmt.Println("error opening file:", errf)
					return false, fmt.Errorf("%w %w", er, err)
				} else {
					linksto = p
				}
			} else {
				h := blake3.New()
				bufFile := bufio.NewReader(file)
				_, errf = io.Copy(h, bufFile)
				if errf != nil {
					file.Close()
					log.Fatal("error:", errf)
				}
				sum = h.Sum(nil)
			}
			file.Close()

			if linksto == in.LinksTo {
			} else if hex.EncodeToString(sum[:]) != in.Blake3Sum {
				return false, fmt.Errorf("blake3 of file %s mismatch", path)
			} else if stat.Mode() != in.FileMode {
				return false, fmt.Errorf("filemode mismatch for file %s", path)
			}
		}
	}
	return true, nil
}
