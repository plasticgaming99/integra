package integrity

import "os"

type IntgFile struct {
	Filepath  string
	FileUid   uint64
	FileGid   uint64
	FileMode  os.FileMode
	Blake3Sum string
}

type Integrity struct {
	ParentPkg string
	Files     []IntgFile
}
