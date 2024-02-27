package entity

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFileGetRequest_Validate(t *testing.T) {
	testCases := []struct {
		name     string
		request  FileGetRequest
		expected error
	}{
		{
			name:     "EmptyFilePath",
			request:  FileGetRequest{FilePath: ""},
			expected: ErrEmptyFilePath,
		},
		{
			name:     "FilePathIsRoot",
			request:  FileGetRequest{FilePath: "/example.txt"},
			expected: ErrForwardSlash,
		},
		{
			name:     "EmptyDirectory",
			request:  FileGetRequest{FilePath: "example.txt"},
			expected: ErrEmptyFilePath,
		},
		{
			name:     "EmptyFileName",
			request:  FileGetRequest{FilePath: "/example/"},
			expected: ErrEmptyFileName,
		},
		{
			name:     "Success",
			request:  FileGetRequest{FilePath: "/example/example.txt"},
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.request.Validate()
			if !assert.Equal(t, tc.expected, err) {
				t.Errorf("expected %v, got %v", tc.expected, err)
			}
		})
	}
}

func TestFileGetResponse_Validate(t *testing.T) {
	testCases := []struct {
		name     string
		filePath string
		fileSize int64
		expected error
	}{
		{
			name:     "EmptyFilePath",
			filePath: "",
			fileSize: 10,
			expected: ErrEmptyFilePath,
		},
		{
			name:     "EmptyFileSize",
			filePath: "test.txt",
			fileSize: 0,
			expected: ErrEmptyContent,
		},
		{
			name:     "Success",
			filePath: "test.txt",
			fileSize: 10,
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp := FileGetResponse{FilePath: tc.filePath, FileSize: tc.fileSize}
			err := resp.Validate()
			if !assert.Equal(t, tc.expected, err) {
				t.Errorf("expected %v, got %v", tc.expected, err)
			}
		})
	}
}
