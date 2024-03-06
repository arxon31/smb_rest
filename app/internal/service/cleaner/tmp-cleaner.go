package cleaner

import (
	"context"
	"git.spbec-mining.ru/arxon31/sambaMW/pkg/logger/sl"
	"log/slog"
	"os"
	"time"
)

type Cleaner struct {
	logger                        *slog.Logger
	tmpDirectoryPath, tmpFilePath string
	timeOffset                    time.Duration
}

func New(logger *slog.Logger, tmpDirectoryPath, tmpFilePath string, timeOffset time.Duration) *Cleaner {
	return &Cleaner{
		logger:           logger,
		tmpDirectoryPath: tmpDirectoryPath,
		tmpFilePath:      tmpFilePath,
		timeOffset:       timeOffset,
	}
}

func (c *Cleaner) Start(ctx context.Context) {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := c.cleanDirs()
			if err != nil {
				c.logger.Error("failed to clean dirs", sl.Err(err))
			}
			err = c.cleanFiles()
			if err != nil {
				c.logger.Error("failed to clean files", sl.Err(err))
			}
		}
	}
}

func (c *Cleaner) cleanDirs() error {
	entries, err := os.ReadDir(c.tmpDirectoryPath)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return err
		}
		dirLifeTime := time.Now().Sub(info.ModTime())
		if dirLifeTime > c.timeOffset {
			err := os.RemoveAll(entry.Name())
			if err != nil {
				return err
			}
			c.logger.Info("removed directory", slog.String("path", entry.Name()))
		}
	}
	return nil
}

func (c *Cleaner) cleanFiles() error {
	entries, err := os.ReadDir(c.tmpFilePath)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return err
		}
		fileLifeTime := time.Now().Sub(info.ModTime())
		if fileLifeTime > c.timeOffset {
			err := os.Remove(entry.Name())
			if err != nil {
				return err
			}
			c.logger.Info("removed file", slog.String("path", entry.Name()))
		}
	}
	return nil
}
