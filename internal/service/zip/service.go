package zip

import (
	"context"
)

type ZipService struct {
}

func (z *ZipService) Zip(ctx context.Context, dirs []string) (zipPath string, err error) {
	return dirs[0], nil
}
