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
			name:        "Empty DirPath",
			req:         DirListRequest{DirPath: ""},
			expectError: ErrEmptyFilePath,
		},
		{
			name:        "DirPath equal to /",
			req:         DirListRequest{DirPath: "/"},
			expectError: ErrRootPath,
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
