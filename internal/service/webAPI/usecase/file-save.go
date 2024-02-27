package usecase

import (
	"context"
	"git.spbec-mining.ru/arxon31/sambaMW/internal/entity"
	"git.spbec-mining.ru/arxon31/sambaMW/pkg/logger/sl"
	"log/slog"
	"path"
	"strings"
)

type FileSaver interface {
	PutFile(ctx context.Context, path string, content []byte) (createdFilePath string, err error)
}

type CacheUpdater interface {
	DeleteEmptyDir(ctx context.Context, dir string) error
}

type FileSaveUsecase struct {
	saver FileSaver
	cache CacheUpdater
	l     *slog.Logger
}

func NewFileSaveUsecase(saver FileSaver, cache CacheUpdater, l *slog.Logger) FileSaveUsecase {
	return FileSaveUsecase{
		saver: saver,
		cache: cache,
		l:     l,
	}
}

func (uc FileSaveUsecase) SaveFile(ctx context.Context, request entity.FileSaveRequest) (response entity.FileSaveResponse, err error) {
	const op = "usecase.FileSaveUsecase.SaveFile()"
	logger := uc.l.With(slog.String("operation", op))

	// Отдаем сохранение файла в слой работы с файлами SAMBA, получаем путь к сохраненному файлу
	logger.Debug("trying to save file", slog.String("path", request.FilePath))
	createdPath, err := uc.saver.PutFile(ctx, request.FilePath, request.Content)
	if err != nil {
		logger.Error("failed to save file", sl.Err(err))
		return entity.FileSaveResponse{}, err
	}

	// Обновляем кэш на предмет нового непустого пути если он был пуст
	logger.Debug("updating cache", slog.String("path", createdPath))

	dir, _ := path.Split(createdPath)
	dir = strings.TrimSuffix(dir, "/")
	err = uc.cache.DeleteEmptyDir(ctx, dir)
	if err != nil {
		logger.Error("failed to update cache", sl.Err(err))
		return entity.FileSaveResponse{}, err
	}

	// Собираем ответ и валидируем его
	resp := entity.FileSaveResponse{FilePath: createdPath}
	logger.Debug("validating saved file", slog.String("path", createdPath))
	err = resp.Validate(request.FilePath)
	if err != nil {
		logger.Error("validation error", sl.Err(err))
		return entity.FileSaveResponse{}, err
	}

	logger.Debug("response", slog.Any("response", resp))
	logger.Info("successfully saved file", slog.String("path", createdPath))

	return resp, nil

}
