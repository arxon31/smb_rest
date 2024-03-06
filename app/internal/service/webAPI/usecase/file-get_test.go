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

var validResponse = entity.FileGetResponse{FilePath: "/tmp/test.txt", FileSize: 1}
var invalidResponse = entity.FileGetResponse{FilePath: "", FileSize: 0}

func TestFileGet(t *testing.T) {
	type fields struct {
		downloader *usecase.MockFileDownloader
	}

	type args struct {
		request entity.FileGetRequest
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
				downloader: usecase.NewMockFileDownloader(gomock.NewController(t)),
			},
			prepare: func(f *fields) {
				f.downloader.EXPECT().GetFile(gomock.Any(), gomock.Any()).Return(validResponse.FilePath, validResponse.FileSize, nil)
			},
			args: args{
				request: entity.FileGetRequest{FilePath: "/tmp/test.txt"},
			},
			wantErr: false,
		},
		{
			name: "Error with file downloading",
			fields: fields{
				downloader: usecase.NewMockFileDownloader(gomock.NewController(t)),
			},
			prepare: func(f *fields) {
				f.downloader.EXPECT().GetFile(gomock.Any(), gomock.Any()).Return("", int64(0), fmt.Errorf("error with downloading file"))
			},
			args: args{
				request: entity.FileGetRequest{FilePath: "/tmp/test.txt"},
			},
			wantErr: true,
		},
		{
			name: "Validation error",
			fields: fields{
				downloader: usecase.NewMockFileDownloader(gomock.NewController(t)),
			},
			prepare: func(f *fields) {
				f.downloader.EXPECT().GetFile(gomock.Any(), gomock.Any()).Return(invalidResponse.FilePath, invalidResponse.FileSize, nil)
			},
			args: args{
				request: entity.FileGetRequest{FilePath: ""},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FileGetUsecase{
				downloader: tt.fields.downloader,
				l:          slog.Default(),
			}
			tt.prepare(&tt.fields)
			if _, err := f.DownloadFile(context.Background(), tt.args.request); (err != nil) != tt.wantErr {
				t.Errorf("FileGetUsecase.DownloadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
