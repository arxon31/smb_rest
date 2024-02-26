package entity

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPaths(t *testing.T) {
	var pathsTestCases = []struct {
		name     string
		expected []string
		node     FileNode
	}{
		{
			name:     "empty node",
			expected: []string{"/"},
			node:     FileNode{},
		},
		{
			name:     "depth 2",
			expected: []string{"/root/folder1", "/root/folder2/subfolder"},
			node: FileNode{
				Name:  "root",
				IsDir: true,
				Child: []FileNode{
					{Name: "folder1", IsDir: true},
					{Name: "folder2", IsDir: true, Child: []FileNode{{Name: "subfolder", IsDir: true}}},
				},
			},
		},
		{
			name:     "depth 3",
			expected: []string{"/root/folder1", "/root/folder2/subfolder/subsubfolder", "/root/folder2/subfolder/subsubfolder2"},
			node: FileNode{
				Name:  "root",
				IsDir: true,
				Child: []FileNode{
					{Name: "folder1", IsDir: true},
					{Name: "folder2", IsDir: true, Child: []FileNode{{Name: "subfolder", IsDir: true, Child: []FileNode{{Name: "subsubfolder", IsDir: true}, {Name: "subsubfolder2", IsDir: true}}}}},
				},
			},
		},
	}
	for _, testCase := range pathsTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			var paths []string
			testCase.node.uniqueFoldersPaths("", &paths)
			if !assert.Equal(t, testCase.expected, paths) {
				t.Errorf("expected paths %v, got %v", testCase.expected, paths)
			}
		})
	}
}
