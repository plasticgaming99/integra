package integrity

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

const testIntegrity = `/
.PACKAGE uid=0 gid=0 perm=644 blake3sum=990582d1c003d7bca6d4aa162b811b2aff82b90d669fb182581e669f4098c592
/usr
/usr/bin
mksh uid=0 gid=0 perm=755 blake3sum=37a1c00d541a2fac64cb7049f91d16b1951d84eb27bef4579d14c5b34cef3277
/usr/share
/usr/share/man
/usr/share/man/man1
mksh.1 uid=0 gid=0 perm=644 blake3sum=41d996859d95c7e0b7b717a2e154e89575a797afb4d066877c2080b030a94939
`

func TestParser(t *testing.T) {
	correctIntg := Integrity{
		ParentPkg: "",
		Files: []IntgFile{
			{
				Filepath:  "/.PACKAGE",
				FileUid:   0,
				FileGid:   0,
				FileMode:  0644,
				Blake3Sum: "990582d1c003d7bca6d4aa162b811b2aff82b90d669fb182581e669f4098c592",
			},
			{
				Filepath:  "/usr/bin/mksh",
				FileUid:   0,
				FileGid:   0,
				FileMode:  0755,
				Blake3Sum: "37a1c00d541a2fac64cb7049f91d16b1951d84eb27bef4579d14c5b34cef3277",
			},
			{
				Filepath:  "/usr/share/man/man1/mksh.1",
				FileUid:   0,
				FileGid:   0,
				FileMode:  0644,
				Blake3Sum: "41d996859d95c7e0b7b717a2e154e89575a797afb4d066877c2080b030a94939",
			},
		},
	}
	rd := strings.NewReader(testIntegrity)
	intg := Parse(rd)
	fmt.Println("correct one:")
	fmt.Println(correctIntg)
	fmt.Println("parser:")
	fmt.Println(intg)
	if !reflect.DeepEqual(intg, correctIntg) {
		t.Errorf("not parsed correctly")
	}
}
