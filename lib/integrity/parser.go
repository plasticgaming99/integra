package integrity

import (
	"bufio"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	FileTypeFile = iota
	FileTypeDir
	FileTypeSym
)

func Parse(in io.Reader) (intg Integrity) {
	dirname := ""
	sc := bufio.NewScanner(in)
	bef, aft, fou := "", "", false
	for sc.Scan() {
		text := sc.Text()
		intgfile := IntgFile{
			FileType: FileTypeFile,
		}

		// perhaps a file
		// first, split file and info
		index := strings.LastIndexByte(text, '/')
		fname := text[:index-1] // exclude " /" string
		args := text[index+1:]  // split with them

		// in intgrity format, directory has / prefix
		// not in file
		if strings.HasPrefix(fname, "/") {
			intgfile.Filepath = fname
			dirname = fname
		} else {
			intgfile.Filepath = filepath.Join(dirname, fname)
		}

		f := strings.Fields(args)
		// loop for attributes
		for _, key := range f {
			bef, aft, fou = strings.Cut(key, "=")
			if !fou {
				// maybe filetype attribute
				switch bef {
				case "dir":
					intgfile.FileType = FileTypeDir
				case "sym":
					intgfile.FileType = FileTypeSym
				}
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
			case "linksto":
				intgfile.LinksTo = aft
			case "blake3sum":
				intgfile.Blake3Sum = aft
			}

		}
		intg.Files = append(intg.Files, intgfile)

	}
	return
}
