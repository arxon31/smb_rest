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
			expected: []string{""},
			node:     FileNode{},
		},
		{
			name: "depth 2",
			expected: []string{
				"root/folder1",
				"root/folder2/subfolder",
			},
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
			name: "depth 3",
			expected: []string{
				"root/folder1",
				"root/folder2/subfolder/subsubfolder",
				"root/folder2/subfolder/subsubfolder2",
			},
			node: FileNode{
				Name:  "root",
				IsDir: true,
				Child: []FileNode{
					{Name: "folder1", IsDir: true},
					{Name: "folder2", IsDir: true, Child: []FileNode{{Name: "subfolder", IsDir: true, Child: []FileNode{{Name: "subsubfolder", IsDir: true}, {Name: "subsubfolder2", IsDir: true}}}}},
				},
			},
		},
		{
			name: "test case",
			expected: []string{
				"OKIS/YUSHKIN/ELIZAR",
				"OKIS/FOCHKIN/VADIM",
				"OKIS/GERASIMENYUK/DAMIR",
				"OKIS/KOZLYUK/VASILY",
			},
			node: FileNode{
				Name:  "OKIS",
				IsDir: true,
				Child: []FileNode{
					{Name: "YUSHKIN", IsDir: true, Child: []FileNode{{Name: "ELIZAR", IsDir: true}}},
					{Name: "FOCHKIN", IsDir: true, Child: []FileNode{{Name: "VADIM", IsDir: true}}},
					{Name: "GERASIMENYUK", IsDir: true, Child: []FileNode{{Name: "DAMIR", IsDir: true}}},
					{Name: "KOZLYUK", IsDir: true, Child: []FileNode{{Name: "VASILY", IsDir: true}}},
				},
			},
		},
	}
	for _, testCase := range pathsTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			res := testCase.node.Paths()
			if !assert.Equal(t, testCase.expected, res) {
				t.Errorf("expected paths %v, got %v", testCase.expected, res)
			}
		})
	}
}
