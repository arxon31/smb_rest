package entity

import (
	"strings"
)

type FileSaveRequest struct {
	FilePath string
	Content  []byte
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

	if s.FilePath != strings.TrimPrefix(s.FilePath, "/") {
		return ErrForwardSlash
	}

	return nil
}

func (s *FileSaveResponse) Validate(requestedFilePath string) error {
	if s.FilePath != requestedFilePath {
		return ErrPathsIsNotEqual
	}
	return nil
}
