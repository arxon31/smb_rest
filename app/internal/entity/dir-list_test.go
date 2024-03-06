package entity

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDirListRequest_Validate(t *testing.T) {
	// Testing for an empty DirPath

	testCases := []struct {
		name        string
		req         DirListRequest
		expectError error
	}{
		{
			name:        "Root path with recursive",
			req:         DirListRequest{DirPath: "", Recursive: true},
			expectError: ErrRootPath,
		},
		{
			name:        "Root path without recursive",
			req:         DirListRequest{DirPath: "", Recursive: false},
			expectError: nil,
		},
		{
			name:        "Forward Slash",
			req:         DirListRequest{DirPath: "/"},
			expectError: ErrForwardSlash,
		},
		{
			name:        "Valid DirPath",
			req:         DirListRequest{DirPath: "valid/path"},
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
