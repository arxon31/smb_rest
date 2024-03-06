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

func TestListDir(t *testing.T) {
	type fields struct {
		viewer *usecase.MockDirectoriesViewer
	}

	type args struct {
		request  entity.DirListRequest
		response entity.DirListResponse
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
				viewer: usecase.NewMockDirectoriesViewer(gomock.NewController(t)),
			},
			prepare: func(f *fields) {
				f.viewer.EXPECT().ListDir(gomock.Any(), gomock.Any(), gomock.Any()).Return(entity.FileNode{}, nil)
			},
			args: args{
				request: entity.DirListRequest{DirPath: tmpDir, Recursive: true},
				response: entity.DirListResponse{Path: entity.FileNode{
					Name:  tmpDir,
					IsDir: true,
					Child: []entity.FileNode{},
				}},
			},
			wantErr: false,
		},
		{
			name: "Error with listing dirs",
			fields: fields{
				viewer: usecase.NewMockDirectoriesViewer(gomock.NewController(t)),
			},
			prepare: func(f *fields) {
				f.viewer.EXPECT().ListDir(gomock.Any(), gomock.Any(), gomock.Any()).Return(entity.FileNode{}, fmt.Errorf("some error"))
			},
			args: args{
				request: entity.DirListRequest{DirPath: "tmp", Recursive: true},
				response: entity.DirListResponse{Path: entity.FileNode{
					Name:  tmpDir,
					IsDir: true,
					Child: []entity.FileNode{},
				}},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare(&tt.fields)
			u := &DirectoryListUsecase{
				viewer: tt.fields.viewer,
				l:      slog.Default(),
			}
			if _, err := u.ListDir(context.Background(), tt.args.request); (err != nil) != tt.wantErr {
				t.Errorf("DirectoryListUsecase.ListDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
