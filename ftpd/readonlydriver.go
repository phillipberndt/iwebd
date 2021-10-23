package ftpd

import (
	"errors"
	"io"

	"goftp.io/server/v2"
	"goftp.io/server/v2/driver/file"
)

// FTP Server is read-only
var readOnly = errors.New("FTP server is read-only")

// Driver for goftp that disallows mutating access.
type readOnlyDriver struct {
	server.Driver
}
func (d *readOnlyDriver) DeleteDir(ctx *server.Context, path string) error { return readOnly; }
func (d *readOnlyDriver) DeleteFile(ctx *server.Context, path string) error { return readOnly; }
func (d *readOnlyDriver) MakeDir(ctx *server.Context, path string) error { return readOnly; }
func (d *readOnlyDriver) PutFile(ctx *server.Context, destPath string, data io.Reader, offset int64) (int64, error) { return 0, readOnly; }

// Create a new file-system backed read-only driver instance
// for use by goftp.
func newReadOnlyDriver(root string) (*readOnlyDriver, error) {
	child, err := file.NewDriver(root)
	if err != nil {
		return nil, err
	}
	return &readOnlyDriver{
		Driver: child,
	}, nil
}
