package cache_updater

import (
	"context"
	"git.spbec-mining.ru/arxon31/sambaMW/internal/entity"
	"git.spbec-mining.ru/arxon31/sambaMW/pkg/logger/sl"
	"log/slog"
)

type Cache interface {
	SaveDirs(ctx context.Context, dirs []string) error
	GetDirs(ctx context.Context) ([]string, error)
	DeleteEmptyDir(ctx context.Context, dir string) error
}

type DirectoryLister interface {
	ListDir(ctx context.Context, dirPath string, recursive bool) (entity.FileNode, error)
}

type Updater struct {
	cache      Cache
	dirsLister DirectoryLister
	logger     *slog.Logger
}

func NewUpdater(cache Cache, dirsLister DirectoryLister) *Updater {
	return &Updater{
		cache:      cache,
		dirsLister: dirsLister,
	}
}

func (u *Updater) update(ctx context.Context) []string {
	dirs, err := u.cache.GetDirs(ctx)
	if err != nil {
		u.logger.Error("failed to get dirs from cache", sl.Err(err))
		return []string{}
	}

	var updatedDirs []string

	for _, dir := range dirs {
		node, err := u.dirsLister.ListDir(ctx, dir, false)
		if err != nil {
			u.logger.Error("failed to list dir", sl.Err(err))
			continue
		}
		if !node.IsEmpty() {
			err = u.cache.DeleteEmptyDir(ctx, dir)
			if err != nil {
				u.logger.Error("failed to delete dir from cache", sl.Err(err))
			}
		}

		updatedDirs = append(updatedDirs, dir)
	}

	return updatedDirs
}
