package integrity

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func Parse(in io.Reader) (intg Integrity) {
	dirname := ""
	sc := bufio.NewScanner(in)
	for sc.Scan() {
		text := sc.Text()
		intgfile := IntgFile{}
		if text[0:1] == "/" {
			// perhaps root dir
			dirname = text
			continue
		} else {
			// perhaps a file
			f := strings.Fields(text)
			for _, key := range f {
				bef, aft, fou := strings.Cut(key, "=")
				if !fou {
					// yas filename
					intgfile.Filepath = filepath.Join(dirname, bef)
					continue
				}
				switch bef {
				case "perm":
					u, err := strconv.ParseUint(aft, 8, 32)
					if err != nil {
						log.Fatal("error:", err)
					}
					intgfile.FileMode = os.FileMode(u)
				case "uid":
					intgfile.FileUid = 0 //currently
				case "gid":
					intgfile.FileGid = 0 //currently
				case "blake3sum":
					intgfile.Blake3Sum = aft
				}

			}
			fmt.Println("append", intgfile)
			intg.Files = append(intg.Files, intgfile)
		}
	}
	return
}
