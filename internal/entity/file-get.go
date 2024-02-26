package entity

import "path"

type FileGetRequest struct {
	FilePath string `json:"filepath"`
}

func (s *FileGetRequest) Validate() error {
	if s.FilePath == "" {
		return ErrEmptyFilePath
	}

	dir, file := path.Split(s.FilePath)
	if dir == "/" {
		return ErrFilePathIsRoot
	}
	if dir == "" {
		return ErrEmptyFilePath
	}
	if file == "" {
		return ErrEmptyFileName
	}

	return nil

}

type FileGetResponse struct {
	FilePath string `json:"filepath"`
	FileSize int64  `json:"filesize"`
}

func (s *FileGetResponse) Validate() error {
	if s.FilePath == "" {
		return ErrEmptyFilePath
	}
	if s.FileSize <= 0 {
		return ErrEmptyContent
	}
	return nil
}
