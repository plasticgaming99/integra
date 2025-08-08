// concurrent integrity
package integrity

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/zeebo/blake3"
)

type Generator struct {
	strb strings.Builder
}

// not needed if you called NewGenerator
func (g *Generator) Init() {
	g.strb.Grow(102400)
}

func NewGenerator() Generator {
	g := Generator{}
	g.Init()
	return g
}

func (g *Generator) Generate(rootpath string) string {
	g.strb.Reset()
	filepath.WalkDir(rootpath, func(path string, d fs.DirEntry, err error) error {
		var errf error
		var inf os.FileInfo
		inf, errf = os.Stat(path)
		if errf != nil {
			log.Fatal("error opening file")
		}
		stat := inf.Sys().(*syscall.Stat_t)

		if d.Type().IsDir() {
			s, errf := filepath.Rel(rootpath, path)
			if errf != nil {
				log.Fatal("error processing path", errf)
			}
			g.strb.WriteString("/")
			if s != "." {
				g.strb.WriteString(s)
			}
			g.strb.WriteString("\n")
		} else if d.Type().Perm().IsRegular() {
			file, errf := os.Open(path)
			if errf != nil {
				log.Fatal("error opening file")
			}
			h := blake3.New()
			bufFile := bufio.NewReader(file)
			_, errf = io.Copy(h, bufFile)
			if errf != nil {
				log.Fatal("error:", errf)
			}
			sum := h.Sum(nil)
			file.Close()

			g.strb.WriteString(d.Name())
			g.strb.WriteString(" ")
			g.strb.WriteString("uid=")
			g.strb.WriteString(fmt.Sprint(stat.Uid))
			g.strb.WriteString(" ")
			g.strb.WriteString("gid=")
			g.strb.WriteString(fmt.Sprint(stat.Gid))
			g.strb.WriteString(" ")
			g.strb.WriteString("perm=")
			g.strb.WriteString(strconv.FormatUint(uint64(inf.Mode().Perm()), 8))
			g.strb.WriteString(" ")
			g.strb.WriteString("blake3sum=")
			g.strb.WriteString(hex.EncodeToString(sum[:]))
			g.strb.WriteString("\n")
		}
		return errf
	})
	return g.strb.String()
}
