package entity

type DirCreateRequest struct {
	Dirs FileNode `json:"path"`
}

func (d DirCreateRequest) Validate() error {

	if d.Dirs.Name == "" {
		return ErrEmptyFilePath
	}

	return nil
}

type DirCreateResponse struct {
	Dirs FileNode `json:"path"`
}

func (d DirCreateResponse) Validate() error {

	if d.Dirs.Name == "" {
		return ErrEmptyFilePath
	}

	return nil
}
