package cmd

import "os"

const (
	DefaultFileMode      os.FileMode = 0660
	DefaultDirectoryMode os.FileMode = 0755
	DefaultExecutable    os.FileMode = 0750
)
