package fileanalyzer

import (
	"os"
)

type Result struct {
	Path   string
	Reason string
}

type FileInfo struct {
	path string
	hash string
	info os.FileInfo
	err  error
}
