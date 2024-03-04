package usecase

import (
	"context"
	"git.spbec-mining.ru/arxon31/sambaMW/internal/entity"
	usecase "git.spbec-mining.ru/arxon31/sambaMW/internal/service/webAPI/usecase/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"os"
	"testing"
)

const (
	tmpPath   = "./tmp"
	tmpDir    = "tmp"
	tmpSubDir = "subdir"
)

func TestCreateDir(t *testing.T) {
	var validRequest = entity.DirCreateRequest{
		Dirs: entity.FileNode{
			IsDir: true,
			Name:  tmpDir,
			Child: []entity.FileNode{
				{
					IsDir: true,
					Name:  tmpSubDir,
				},
			},
		},
	}
	var validResponse = entity.DirCreateResponse{
		Dirs: entity.FileNode{
			IsDir: true,
			Name:  tmpDir,
			Child: []entity.FileNode{
				{
					IsDir: true,
					Name:  tmpSubDir,
				},
			},
		},
	}

	tests := []struct {
		name     string
		request  entity.DirCreateRequest
		response entity.DirCreateResponse
		err      error
	}{
		{
			name:     "Successful creation",
			request:  validRequest,
			response: validResponse,
			err:      nil,
		},
	}

	// Test when paths are empty
	//t.Run("EmptyPaths", func(t *testing.T) {
	//	// Test logic here
	//})
	//
	//// Test when all directories are successfully created
	//t.Run("AllDirsCreatedSuccessfully", func(t *testing.T) {
	//	// Test logic here
	//})
	//
	//// Test when some directories fail to be created
	//t.Run("SomeDirsFailedToCreate", func(t *testing.T) {
	//	// Test logic here
	//})
	//
	//// Test when cache.SaveDirs returns an error
	//t.Run("CacheSaveDirsError", func(t *testing.T) {
	//	// Test logic here
	//})
	//
	//// Test when response validation fails
	//t.Run("ResponseValidationFailure", func(t *testing.T) {
	//	// Test logic here
	//})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testContext := context.Background()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			creator := usecase.NewMockDirectoryCreator(ctrl)
			cacher := usecase.NewMockEmptyDirsCache(ctrl)
			paths := tt.request.Dirs.Paths()

			creator.EXPECT().CreateDir(testContext, paths[0]).Return(paths[0], nil)
			cacher.EXPECT().SaveDirs(testContext, tt.response.Dirs.Paths()).Return(nil)

			usecase := NewDirectoryCreateUsecase(creator, cacher, slog.New(slog.NewTextHandler(os.Stdout, nil)))

			resp, err := usecase.CreateDir(testContext, tt.request)

			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.response, resp)

		})
	}
}
