package usecase

import (
	"context"
	"git.spbec-mining.ru/arxon31/sambaMW/internal/entity"
	"git.spbec-mining.ru/arxon31/sambaMW/pkg/logger/sl"
	"log/slog"
)

//go:generate mockgen -source=dir-list.go -destination=mocks/dir-list-mock.go -package=mocks
type DirectoriesViewer interface {
	ListDir(ctx context.Context, dirPath string, recursive bool) (entity.FileNode, error)
}

type DirectoryListUsecase struct {
	viewer DirectoriesViewer
	l      *slog.Logger
}

func NewDirectoryListUsecase(viewer DirectoriesViewer, logger *slog.Logger) DirectoryListUsecase {
	return DirectoryListUsecase{
		viewer: viewer,
		l:      logger,
	}
}

func (d DirectoryListUsecase) ListDir(ctx context.Context, request entity.DirListRequest) (response entity.DirListResponse, err error) {
	const op = "usecase.DirectoryListUsecase.ListDir()"
	logger := d.l.With(slog.String("operation", op))

	logger.Debug("trying to list dir", slog.String("path", request.DirPath))

	node, err := d.viewer.ListDir(ctx, request.DirPath, request.Recursive)
	if err != nil {
		logger.Error("failed to list dir", sl.Err(err))
		return entity.DirListResponse{}, err
	}

	logger.Debug("got node", slog.Any("node", node))
	resp := entity.DirListResponse{Path: node}

	err = resp.Validate()
	if err != nil {
		logger.Error("validation error", sl.Err(err))
		return entity.DirListResponse{}, err
	}

	logger.Debug("response", slog.Any("response", resp))
	logger.Info("successfully listed dir", slog.String("path", request.DirPath))

	return resp, nil

}
