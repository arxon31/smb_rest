package entity

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDirCreateRequest_Validate(t *testing.T) {

	testCases := []struct {
		name        string
		req         DirCreateRequest
		expectError error
	}{
		{
			name:        "Empty DirPath",
			req:         DirCreateRequest{Dirs: FileNode{Name: ""}},
			expectError: ErrEmptyFilePath,
		},
		{
			name:        "Valid DirPath",
			req:         DirCreateRequest{Dirs: FileNode{Name: "valid/path"}},
			expectError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.req.Validate()
			if !assert.Equal(t, tc.expectError, err) {
				t.Errorf("expected error %v, got %v", tc.expectError, err)
			}
		})
	}
}

func TestDirCreateResponse_Validate(t *testing.T) {

	testCases := []struct {
		name        string
		res         DirCreateResponse
		expectError error
	}{
		{
			name:        "Empty DirPath",
			res:         DirCreateResponse{Dirs: FileNode{Name: ""}},
			expectError: ErrEmptyFilePath,
		},
		{
			name:        "Valid DirPath",
			res:         DirCreateResponse{Dirs: FileNode{Name: "valid/path"}},
			expectError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.res.Validate()
			if !assert.Equal(t, tc.expectError, err) {
				t.Errorf("expected error %v, got %v", tc.expectError, err)
			}
		})
	}

}
