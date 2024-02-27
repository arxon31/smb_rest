package entity

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFileSaveRequest_Validate(t *testing.T) {
	testCases := []struct {
		name     string
		request  FileSaveRequest
		expected error
	}{
		{
			name:     "EmptyFilePath",
			request:  FileSaveRequest{FilePath: "", Content: []byte("test")},
			expected: ErrEmptyFilePath,
		},
		{
			name:     "EmptyContent",
			request:  FileSaveRequest{FilePath: "/path/to/file", Content: []byte{}},
			expected: ErrEmptyContent,
		},
		{
			name:     "ForwardSlash",
			request:  FileSaveRequest{FilePath: "/", Content: []byte("test")},
			expected: ErrForwardSlash,
		},
		{
			name:     "Success",
			request:  FileSaveRequest{FilePath: "path/to/file", Content: []byte("test")},
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.request.Validate()
			if !assert.Equal(t, tc.expected, err) {
				t.Errorf("expected error %v, got %v", tc.expected, err)
			}
		})
	}
}

func TestFileSaveResponse_Validate(t *testing.T) {
	testCases := []struct {
		name          string
		filePath      string
		requestedPath string
		expectedError error
	}{
		{"success", "path1", "path1", nil},
		{"fail", "path1", "path2", ErrPathsIsNotEqual},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp := FileSaveResponse{FilePath: tc.filePath}
			err := resp.Validate(tc.requestedPath)
			if !assert.Equal(t, tc.expectedError, err) {
				t.Errorf("expected error %v, got %v", tc.expectedError, err)
			}
		})
	}
}
