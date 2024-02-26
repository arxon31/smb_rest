package entity

type FileNode struct {
	Name   string     `json:"name"`
	IsDir  bool       `json:"is_dir"`
	SaveAs string     `json:"save_as,omitempty"`
	Child  []FileNode `json:"child,omitempty"`
}

func (t *FileNode) Paths() []string {
	var paths []string

	t.uniqueFoldersPaths(t.Name, &paths)

	return paths
}

func (t *FileNode) uniqueFoldersPaths(path string, paths *[]string) {
	path += "/" + t.Name
	if len(t.Child) == 0 {
		*paths = append(*paths, path)
		return
	}

	for _, child := range t.Child {

		child.uniqueFoldersPaths(path, paths)
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
