package integrity

import "os"

type IntgFile struct {
	Filepath  string
	FileType  int
	FileUid   uint64
	FileGid   uint64
	FileMode  os.FileMode
	LinksTo   string
	Blake3Sum string
}

type Integrity struct {
	ParentPkg string
	Files     []IntgFile
}
