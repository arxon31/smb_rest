package usecase

import (
	"context"
	"fmt"
	"git.spbec-mining.ru/arxon31/sambaMW/internal/entity"
	usecase "git.spbec-mining.ru/arxon31/sambaMW/internal/service/webAPI/usecase/mocks"
	"github.com/golang/mock/gomock"
	"log/slog"
	"testing"
)

const (
	tmpDir    = "tmp"
	tmpSubDir = "subdir"
)

var request = entity.DirCreateRequest{
	Dirs: entity.FileNode{
		Name:  tmpDir,
		IsDir: true,
		Child: []entity.FileNode{
			{
				Name:  tmpSubDir,
				IsDir: true,
			},
		},
	},
}

var response = entity.DirCreateResponse{
	Dirs: entity.FileNode{
		Name:  tmpDir,
		IsDir: true,
		Child: []entity.FileNode{
			{
				Name:  tmpSubDir,
				IsDir: true,
			},
		},
	},
}

func TestCreateDir(t *testing.T) {

	type fields struct {
		creator *usecase.MockDirectoryCreator
		cache   *usecase.MockEmptyDirsCache
	}

	type args struct {
		request  entity.DirCreateRequest
		response entity.DirCreateResponse
	}

	tests := []struct {
		fields  fields
		name    string
		prepare func(f *fields)
		args    args
		wantErr bool
	}{
		{
			name: "Success",
			fields: fields{
				creator: usecase.NewMockDirectoryCreator(gomock.NewController(t)),
				cache:   usecase.NewMockEmptyDirsCache(gomock.NewController(t)),
			},
			prepare: func(f *fields) {
				f.creator.EXPECT().CreateDir(gomock.Any(), gomock.Any()).Return([]string{}, nil)
				f.cache.EXPECT().SaveDirs(gomock.Any(), gomock.Any()).Return(nil)
			}, args: args{
				request:  request,
				response: response,
			},
			wantErr: false,
		},
		{
			name: "Error with creating dirs",
			fields: fields{
				creator: usecase.NewMockDirectoryCreator(gomock.NewController(t)),
				cache:   usecase.NewMockEmptyDirsCache(gomock.NewController(t)),
			},
			prepare: func(f *fields) {
				f.creator.EXPECT().CreateDir(gomock.Any(), gomock.Any()).Return([]string{}, fmt.Errorf("some error"))
			},
			args: args{
				request:  request,
				response: response,
			},
			wantErr: true,
		},
		{
			name: "Error with saving dirs to cache",
			fields: fields{
				creator: usecase.NewMockDirectoryCreator(gomock.NewController(t)),
				cache:   usecase.NewMockEmptyDirsCache(gomock.NewController(t)),
			},
			prepare: func(f *fields) {
				f.creator.EXPECT().CreateDir(gomock.Any(), gomock.Any()).Return([]string{}, nil)
				f.cache.EXPECT().SaveDirs(gomock.Any(), gomock.Any()).Return(fmt.Errorf("some error"))
			},
			args: args{
				request:  request,
				response: response,
			},
			wantErr: true,
		},
		{
			name: "Error with validation response",
			fields: fields{
				creator: usecase.NewMockDirectoryCreator(gomock.NewController(t)),
				cache:   usecase.NewMockEmptyDirsCache(gomock.NewController(t)),
			},
			prepare: func(f *fields) {
				f.creator.EXPECT().CreateDir(gomock.Any(), gomock.Any()).Return([]string{}, nil)
				f.cache.EXPECT().SaveDirs(gomock.Any(), gomock.Any()).Return(nil)
			},
			args: args{
				request: entity.DirCreateRequest{Dirs: entity.FileNode{
					Name:  "",
					IsDir: true,
					Child: []entity.FileNode{},
				}},
				response: entity.DirCreateResponse{Dirs: entity.FileNode{
					Name:  "",
					IsDir: true,
					Child: []entity.FileNode{},
				}},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testContext := context.Background()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			tt.prepare(&tt.fields)

			uCase := DirectoryCreateUsecase{
				creator: tt.fields.creator,
				cache:   tt.fields.cache,
				l:       slog.Default(),
			}

			_, err := uCase.DirectoryCreate(testContext, tt.args.request)

			if (err != nil) != tt.wantErr {
				t.Errorf("DirectoryCreate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
