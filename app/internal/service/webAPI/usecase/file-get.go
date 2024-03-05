package usecase

import (
	"context"
	"git.spbec-mining.ru/arxon31/sambaMW/internal/entity"
	"git.spbec-mining.ru/arxon31/sambaMW/pkg/logger/sl"
	"log/slog"
)

//go:generate mockgen -source=file-get.go -destination=mocks/file-get-mock.go -package=mocks
type FileDownloader interface {
	GetFile(ctx context.Context, path string) (filename string, filesize int64, err error)
}

type FileGetUsecase struct {
	downloader FileDownloader
	l          *slog.Logger
}

func NewFileGetUsecase(downloader FileDownloader, logger *slog.Logger) FileGetUsecase {
	return FileGetUsecase{
		downloader: downloader,
		l:          logger,
	}
}

func (f FileGetUsecase) DownloadFile(ctx context.Context, request entity.FileGetRequest) (response entity.FileGetResponse, err error) {
	const op = "usecase.FileGetUsecase.DownloadFile()"
	logger := f.l.With(slog.String("operation", op))

	logger.Debug("trying to get file", slog.String("path", request.FilePath))
	filePath, filesize, err := f.downloader.GetFile(ctx, request.FilePath)
	if err != nil {
		logger.Error("failed to get file", sl.Err(err))
		return entity.FileGetResponse{}, err
	}

	logger.Debug("got file", slog.String("filePath", filePath), slog.Int64("filesize", filesize))
	resp := entity.FileGetResponse{FilePath: filePath, FileSize: filesize}
	logger.Debug("validating gotten file", slog.String("filePath", filePath))
	err = resp.Validate()
	if err != nil {
		logger.Error("validation error", sl.Err(err))
		return entity.FileGetResponse{}, err
	}

	logger.Debug("response", slog.Any("response", resp))
	logger.Info("successfully got file", slog.String("filePath", filePath), slog.Int64("filesize", filesize))

	return resp, nil

}
