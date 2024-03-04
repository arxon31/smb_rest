package samba

import (
	"context"
	"git.spbec-mining.ru/arxon31/sambaMW/internal/entity"
	"git.spbec-mining.ru/arxon31/sambaMW/pkg/samba/sessions"
	"github.com/hirochachacha/go-smb2"
	"io"
	"log/slog"
	"os"
	"path"
	"strings"
)

const (
	permissions = 0777
)

type smbSessionManager interface {
	GetSession() (sessions.Session, error)
	ReleaseSession(session sessions.Session)
}

type Client struct {
	sm                 smbSessionManager
	logger             *slog.Logger
	shareName          string
	tmpDirectoriesPath string
	tmpFilesPath       string
}

func New(logger *slog.Logger, host, port string, user, password, shareName string, poolSize int, tmpDirectoriesPath, tmpFilesPath string) (*Client, error) {

	cm, err := sessions.NewSessionManager(logger, host, port, user, password, poolSize)
	if err != nil {
		return &Client{}, err
	}

	err = os.MkdirAll(tmpDirectoriesPath, permissions)
	if err != nil {
		return &Client{}, err
	}
	err = os.MkdirAll(tmpFilesPath, permissions)
	if err != nil {
		return &Client{}, err
	}

	return &Client{sm: cm, shareName: shareName, logger: logger, tmpDirectoriesPath: tmpDirectoriesPath, tmpFilesPath: tmpFilesPath}, nil

}
func (c *Client) ListDir(ctx context.Context, dirPath string, recursive bool) (entity.FileNode, error) {
	session, err := c.sm.GetSession()
	if err != nil {
		return entity.FileNode{}, err
	}
	defer c.sm.ReleaseSession(session)

	fs, err := session.Mount(c.shareName)
	if err != nil {
		return entity.FileNode{}, err
	}
	defer fs.Umount()

	if !recursive {
		return c.makeFileTree(ctx, fs, dirPath)
	}

	return c.makeFileTreeAll(ctx, fs, dirPath)

}

func (c *Client) CreateDir(ctx context.Context, dirPath string) (createdDirPath string, err error) {
	session, err := c.sm.GetSession()
	if err != nil {
		return "", err
	}
	defer c.sm.ReleaseSession(session)

	fs, err := session.Mount(c.shareName)
	if err != nil {
		return "", err
	}
	defer fs.Umount()

	for _, dir := range dirs {
		err = fs.WithContext(ctx).MkdirAll(dirPath, permissions)
		if err != nil {
			return "", err
		}
	}

	return dirPath, nil
}

func (c *Client) GetDirectory(ctx context.Context, dirPath string, saveAs string) (tempDirPath string, err error) {
	session, err := c.sm.GetSession()
	if err != nil {
		return "", err
	}
	defer c.sm.ReleaseSession(session)

	fs, err := session.Mount(c.shareName)
	if err != nil {
		return "", err
	}
	defer fs.Umount()

	dstPath := path.Join(c.tmpDirectoriesPath, saveAs)

	err = c.copyDirAll(ctx, fs, dirPath, dstPath)
	if err != nil {
		return "", err
	}

	return dstPath, nil
}

func (c *Client) GetFile(ctx context.Context, filePath string) (createdPath string, size int64, err error) {
	session, err := c.sm.GetSession()
	if err != nil {
		return "", 0, err
	}
	defer c.sm.ReleaseSession(session)

	fs, err := session.Mount(c.shareName)
	if err != nil {
		return "", 0, err
	}
	defer fs.Umount()

	_, err = fs.Stat(filePath)
	if err != nil || os.IsNotExist(err) {
		return "", 0, os.ErrNotExist
	}

	_, fileName := path.Split(filePath)

	dstFilePath := path.Join(c.tmpFilesPath, fileName)

	fileSize, err := c.copyFile(ctx, fs, filePath, dstFilePath)
	if err != nil {
		return "", 0, err
	}

	return dstFilePath, fileSize, nil
}

func (c *Client) PutFile(ctx context.Context, path string, content []byte) (createdFilePath string, err error) {
	session, err := c.sm.GetSession()
	if err != nil {
		return "", err
	}
	defer c.sm.ReleaseSession(session)

	fs, err := session.Mount(c.shareName)
	if err != nil {
		return "", err
	}
	defer fs.Umount()

	err = fs.WithContext(ctx).WriteFile(path, content, permissions)
	if err != nil {
		return "", err
	}

	return path, nil

}

func (c *Client) makeFileTree(ctx context.Context, fs *smb2.Share, dirPath string) (entity.FileNode, error) {
	p := strings.Split(dirPath, "/")

	info, err := fs.WithContext(ctx).Lstat(dirPath)
	if err != nil {
		return entity.FileNode{}, err
	}

	files := entity.FileNode{
		Name:  p[len(p)-1],
		IsDir: info.IsDir(),
	}
	if info.IsDir() {
		entries, err := fs.WithContext(ctx).ReadDir(dirPath)
		if err != nil {
			return entity.FileNode{}, err
		}
		for _, entry := range entries {
			files.Child = append(files.Child, entity.FileNode{
				Name:  entry.Name(),
				IsDir: entry.IsDir(),
			})
		}
	}
	return files, nil
}

func (c *Client) makeFileTreeAll(ctx context.Context, fs *smb2.Share, dirPath string) (entity.FileNode, error) {
	p := strings.Split(dirPath, "/")
	info, err := fs.WithContext(ctx).Lstat(dirPath)
	if err != nil {
		return entity.FileNode{}, err
	}

	files := entity.FileNode{
		Name:  p[len(p)-1],
		IsDir: info.IsDir(),
	}

	if info.IsDir() {
		entries, err := fs.WithContext(ctx).ReadDir(dirPath)
		if err != nil {
			return entity.FileNode{}, err
		}
		for _, entry := range entries {
			if entry.IsDir() {

				childPath := path.Join(dirPath, entry.Name())
				node, err := c.makeFileTreeAll(ctx, fs, childPath)
				if err != nil {
					return entity.FileNode{}, err
				}
				files.Child = append(files.Child, node)
			} else {
				files.Child = append(files.Child, entity.FileNode{
					Name:  entry.Name(),
					IsDir: entry.IsDir(),
				})
			}
		}
	}

	return files, nil

}

func (c *Client) copyFile(ctx context.Context, samba *smb2.Share, src, dst string) (int64, error) {
	srcFile, err := samba.WithContext(ctx).Open(src)
	if err != nil {
		return 0, err
	}
	defer srcFile.Close()

	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_RDWR, permissions)
	if err != nil {
		return 0, err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return 0, err
	}

	info, err := srcFile.Stat()
	if err != nil {
		return 0, err
	}

	return info.Size(), nil
}

func (c *Client) copyDirAll(ctx context.Context, samba *smb2.Share, src, dst string) error {
	info, err := samba.Stat(src)
	if err != nil {
		return err
	}

	err = os.MkdirAll(dst, info.Mode())
	if err != nil {
		return err
	}

	entries, err := samba.WithContext(ctx).ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := path.Join(src, entry.Name())
		dstPath := path.Join(dst, entry.Name())
		if entry.IsDir() {
			err = c.copyDirAll(ctx, samba, srcPath, dstPath)
			if err != nil {
				return err
			}
		} else {
			_, err = c.copyFile(ctx, samba, srcPath, dstPath)
			if err != nil {
				return err
			}
		}

	}
	return nil
}
