package entity

import "errors"

var (
	ErrEmptyFilePath   = errors.New("empty file path")
	ErrEmptyContent    = errors.New("empty content")
	ErrForwardSlash    = errors.New("forward slash in path")
	ErrEmptyFileName   = errors.New("empty file name")
	ErrPathsIsNotEqual = errors.New("paths are not equal")
	ErrRootPath        = errors.New("root path detected")
	ErrInvalidUser     = errors.New("invalid user")
)
