package usecase

import (
	"context"
	"git.spbec-mining.ru/arxon31/sambaMW/internal/entity"
	"git.spbec-mining.ru/arxon31/sambaMW/pkg/logger/sl"
	"log/slog"
)

type DirectoryDownloader interface {
	GetDirectory(ctx context.Context, dirPath string, saveAs string) (tempDirPath string, err error)
}

type Zipper interface {
	Zip(ctx context.Context, dirs []string) (zipPath string, err error)
}

type DirectoryGetUsecase struct {
	downloader DirectoryDownloader
	zipper     Zipper
	l          *slog.Logger
}

func NewDirectoryGetUsecase(downloader DirectoryDownloader, zipper Zipper, logger *slog.Logger) DirectoryGetUsecase {
	return DirectoryGetUsecase{
		downloader: downloader,
		zipper:     zipper,
		l:          logger,
	}
}

// TODO: think about this method
func (d DirectoryGetUsecase) GetDirectory(ctx context.Context, request entity.DirGetRequest) (response entity.DirGetResponse, err error) {
	const op = "usecase.DirectoryGetUsecase.GetDirectory()"
	logger := d.l.With(slog.String("operation", op))

	var createdDirs []string

	for saveAs, node := range request.Dirs {
		paths := node.Paths()
		for _, path := range paths {
			tmpCreatedDir, err := d.downloader.GetDirectory(ctx, path, saveAs)
			if err != nil {
				logger.Error("failed to get dir", sl.Err(err))
				return entity.DirGetResponse{}, err
			}
			createdDirs = append(createdDirs, tmpCreatedDir)
		}
	}

	logger.Debug("starting to zip directory", slog.String("path", ""))
	zippedDirs, err := d.zipper.Zip(ctx, createdDirs)
	if err != nil {
		logger.Error("failed to zip directory", sl.Err(err))
		return entity.DirGetResponse{}, err
	}

	resp := entity.DirGetResponse{DirPath: zippedDirs}
	err = resp.Validate()
	if err != nil {
		logger.Error("failed to validate response", sl.Err(err))
		return entity.DirGetResponse{}, err
	}

	logger.Debug("response", slog.Any("response", resp))
	logger.Info("successfully got directory", slog.String("path", zippedDirs))

	return resp, err

}
