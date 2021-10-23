package webdavd

import (
	"errors"
	"context"
	"os"

	"golang.org/x/net/webdav"
)

var forbidden = errors.New("Forbidden")

type readOnlyDir struct {
	webdav.Dir
}

func (d readOnlyDir) Mkdir(c context.Context, name string, perm os.FileMode) error {
	return forbidden
}

func (d readOnlyDir) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	if flag & (os.O_RDWR | os.O_WRONLY) != 0 {
		return nil, forbidden
	}
	return d.Dir.OpenFile(ctx, name, flag, perm)
}

func (d readOnlyDir) RemoveAll(ctx context.Context, name string) error {
	return forbidden
}

func (d readOnlyDir) Rename(ctx context.Context, oldName, newName string) error {
	return forbidden
}
