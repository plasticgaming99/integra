package cmd

import (
	"fmt"
	"io"
)

const HelpMsg = `Integra Help
UpperCase Options:
  I nstall : install package
  S ync    : synchronize database
  U pgrade : upgrade package
  R emove  : remove package

Normal Options:
  --install                 : install package
  --sync                    : synchronize database
  --search                  : search package
  --upgrade                 : upgrade package
  --remove                  : remove package
  --yes                     : continue automatically
  --dbg                     : enable debug printing
  --verbose                 : output more message
  --quiet                   : output less message
  --override-root=[rootdir] : override root directory for bootstrapping`

func CommandHelp(w io.Writer) (err error) {
	_, err = fmt.Fprintln(w, HelpMsg)
	return
}
