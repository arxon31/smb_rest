package entity

type DirListRequest struct {
	DirPath   string `json:"path"`
	Recursive bool   `json:"recursive"`
}

func (d DirListRequest) Validate() error {

	if d.DirPath == "" {
		return ErrEmptyFilePath
	}

	if d.DirPath == "/" {
		return ErrRootPath
	}

	return nil
}

type DirListResponse struct {
	Path FileNode `json:"path"`
}

func (d DirListResponse) Validate() error {

	return nil
}
