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

func TestFileSave(t *testing.T) {
	type fields struct {
		saver *usecase.MockFileSaver
		cache *usecase.MockCacheUpdater
	}

	type args struct {
		request entity.FileSaveRequest
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
				saver: usecase.NewMockFileSaver(gomock.NewController(t)),
				cache: usecase.NewMockCacheUpdater(gomock.NewController(t)),
			},
			prepare: func(f *fields) {
				f.saver.EXPECT().PutFile(gomock.Any(), gomock.Any(), gomock.Any()).Return("/tmp/test.txt", nil)
				f.cache.EXPECT().DeleteEmptyDir(gomock.Any(), gomock.Any()).Return(nil)
			},
			args: args{
				request: entity.FileSaveRequest{FilePath: "/tmp/test.txt", Content: []byte("test")},
			},
			wantErr: false,
		},
		{
			name: "Error with saving file",
			fields: fields{
				saver: usecase.NewMockFileSaver(gomock.NewController(t)),
				cache: usecase.NewMockCacheUpdater(gomock.NewController(t)),
			},
			prepare: func(f *fields) {
				f.saver.EXPECT().PutFile(gomock.Any(), gomock.Any(), gomock.Any()).Return("", fmt.Errorf("error with saving file"))
			},
			args: args{
				request: entity.FileSaveRequest{FilePath: "/tmp/test.txt", Content: []byte("test")},
			},
			wantErr: true,
		},
		{
			name: "Error with deleting empty dir",
			fields: fields{
				saver: usecase.NewMockFileSaver(gomock.NewController(t)),
				cache: usecase.NewMockCacheUpdater(gomock.NewController(t)),
			},
			prepare: func(f *fields) {
				f.saver.EXPECT().PutFile(gomock.Any(), gomock.Any(), gomock.Any()).Return("/tmp/test.txt", nil)
				f.cache.EXPECT().DeleteEmptyDir(gomock.Any(), gomock.Any()).Return(fmt.Errorf("error with deleting empty dir"))
			},
			args: args{
				request: entity.FileSaveRequest{FilePath: "/tmp/test.txt", Content: []byte("test")},
			},
			wantErr: true,
		},
		{
			name: "Error with validation",
			fields: fields{
				saver: usecase.NewMockFileSaver(gomock.NewController(t)),
				cache: usecase.NewMockCacheUpdater(gomock.NewController(t)),
			},
			prepare: func(f *fields) {
				f.saver.EXPECT().PutFile(gomock.Any(), gomock.Any(), gomock.Any()).Return("/tmp/test.txt", nil)
				f.cache.EXPECT().DeleteEmptyDir(gomock.Any(), gomock.Any()).Return(nil)
			},
			args: args{
				request: entity.FileSaveRequest{FilePath: "/tmp/test/test.txt", Content: []byte("test")},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := FileSaveUsecase{
				saver: tt.fields.saver,
				cache: tt.fields.cache,
				l:     slog.Default(),
			}
			tt.prepare(&tt.fields)
			if _, err := f.SaveFile(context.Background(), tt.args.request); (err != nil) != tt.wantErr {
				t.Errorf("FileSaver.PutFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
