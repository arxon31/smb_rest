package entity

import (
	"path"
	"strings"
)

type FileNode struct {
	Name   string     `json:"name"`
	IsDir  bool       `json:"is_dir"`
	SaveAs string     `json:"save_as,omitempty"`
	Child  []FileNode `json:"child,omitempty"`
}

func (t *FileNode) Paths() []string {
	var paths []string

	t.uniqueFoldersPaths("", &paths)

	for i, path := range paths {
		paths[i] = strings.TrimPrefix(path, "/")
	}

	return paths
}

func (t *FileNode) uniqueFoldersPaths(p string, paths *[]string) {
	
	p = path.Join(p, t.Name)
	if len(t.Child) == 0 {
		*paths = append(*paths, p)
		return
	}

	for _, child := range t.Child {

		child.uniqueFoldersPaths(p, paths)
	}
}

func (t *FileNode) IsEmpty() bool {
	if len(t.Child) == 0 {
		return true
	}

	for _, child := range t.Child {
		if !child.IsDir {
			return false
		}
	}

	return true

}
