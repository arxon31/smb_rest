package entity

import "strings"

type DirListRequest struct {
	DirPath   string `json:"path"`
	Recursive bool   `json:"recursive"`
}

func (d DirListRequest) Validate() error {

	if d.DirPath == "" && d.Recursive {
		return ErrRootPath
	}

	if d.DirPath != strings.TrimPrefix(d.DirPath, "/") {
		return ErrForwardSlash
	}

	return nil
}

type DirListResponse struct {
	Path FileNode `json:"path"`
}

func (d DirListResponse) Validate() error {

	return nil
}
