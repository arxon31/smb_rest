package usecase

import (
	"context"
	"git.spbec-mining.ru/arxon31/sambaMW/internal/entity"
	"git.spbec-mining.ru/arxon31/sambaMW/pkg/logger/sl"
	"log/slog"
	"slices"
)

//go:generate mockgen -source=dir-create.go -destination=mocks/dir-create-mock.go -package=mocks
type DirectoryCreator interface {
	CreateDir(ctx context.Context, dirs []string) (createdDirs []string, err error)
}

type EmptyDirsCache interface {
	SaveDirs(ctx context.Context, dirs []string) error
}

type DirectoryCreateUsecase struct {
	creator DirectoryCreator
	cache   EmptyDirsCache
	l       *slog.Logger
}

func NewDirectoryCreateUsecase(creator DirectoryCreator, cache EmptyDirsCache, logger *slog.Logger) DirectoryCreateUsecase {
	return DirectoryCreateUsecase{
		creator: creator,
		cache:   cache,
		l:       logger,
	}
}

func (d DirectoryCreateUsecase) CreateDir(ctx context.Context, request entity.DirCreateRequest) (response entity.DirCreateResponse, err error) {
	const op = "usecase.DirectoryCreateUsecase.CreateDir()"
	logger := d.l.With(slog.String("operation", op))

	logger.Debug("trying to extract unique paths from FileNode", slog.Any("paths", request.Dirs))
	paths := request.Dirs.Paths()

	logger.Debug("trying to create dirs", slog.Any("paths", paths))
	createdDirs, err := d.creator.CreateDir(ctx, paths)
	if err != nil {
		logger.Error("failed to create dir", sl.Err(err))
		return entity.DirCreateResponse{}, err
	}

	logger.Debug("created dirs",
		slog.Bool("all-success", slices.Equal(createdDirs, paths)))

	err = d.cache.SaveDirs(ctx, createdDirs)
	if err != nil {
		logger.Error("failed to save created dirs", sl.Err(err))
		return entity.DirCreateResponse{}, err
	}

	resp := entity.DirCreateResponse{Dirs: request.Dirs}
	logger.Debug("validating created dirs", slog.Any("paths", createdDirs))
	err = resp.Validate()
	if err != nil {
		logger.Error("failed to validate response", sl.Err(err))
		return entity.DirCreateResponse{}, err
	}

	logger.Debug("response", slog.Any("response", createdDirs))
	logger.Info("successfully created dirs", slog.Any("paths", createdDirs))

	return resp, nil

}
