package entity

import "path"

type FileSaveRequest struct {
	FilePath string `json:"filepath"`
	Content  []byte `json:"content"`
}

type FileSaveResponse struct {
	FilePath string `json:"filepath"`
}

func (s *FileSaveRequest) Validate() error {
	if s.FilePath == "" {
		return ErrEmptyFilePath
	}
	if s.Content == nil || len(s.Content) == 0 {
		return ErrEmptyContent
	}

	dir, file := path.Split(s.FilePath)
	if dir == "/" {
		return ErrForwardSlash
	}
	if dir == "" {
		return ErrEmptyFilePath
	}
	if file == "" {
		return ErrEmptyFileName
	}

	return nil
}

func (s *FileSaveResponse) Validate(requestedFilePath string) error {
	if s.FilePath != requestedFilePath {
		return ErrPathsIsNotEqual
	}
	return nil
}
