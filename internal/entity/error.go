package entity

import "errors"

var (
	ErrEmptyFilePath   = errors.New("empty file path")
	ErrEmptyContent    = errors.New("empty content")
	ErrFilePathIsRoot  = errors.New("file path is root")
	ErrEmptyFileName   = errors.New("empty file name")
	ErrPathsIsNotEqual = errors.New("paths are not equal")
	ErrRootPath        = errors.New("root path")
)
