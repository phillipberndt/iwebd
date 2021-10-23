package ftpd

import (
	"errors"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"syscall"
)

// Operation not supported.
var unsupported = errors.New("Unsupported operation")

// Implements the permission interface for goftp.io/v2
// for basic POSIX file systems.
type fsPerm struct {
	root string
}

// Return a new file system permission interface instance for the given root
// directory.
func newFsPerm(root string) *fsPerm {
	return &fsPerm{
		root,
	}
}
func (f *fsPerm) path(name string) string {
	return filepath.Join(f.root, name)
}
func (f *fsPerm) GetOwner(name string) (string, error) {
	finfo, err := os.Stat(f.path(name))
	if err != nil {
		return "", err
	}
	stat, ok := finfo.Sys().(*syscall.Stat_t)
	if !ok {
		return "", unsupported
	}
	uid := strconv.FormatInt(int64(stat.Uid), 10)
	user, err := user.LookupId(uid)
	if err != nil {
		return uid, nil
	}
	return user.Name, nil

}
func (f *fsPerm) GetGroup(name string) (string, error) {
	finfo, err := os.Stat(f.path(name))
	if err != nil {
		return "", err
	}
	stat, ok := finfo.Sys().(*syscall.Stat_t)
	if !ok {
		return "", unsupported
	}
	gid := strconv.FormatInt(int64(stat.Gid), 10)
	group, err := user.LookupGroupId(gid)
	if err != nil {
		return gid, nil
	}
	return group.Name, nil
}
func (f *fsPerm) GetMode(name string) (os.FileMode, error) {
	finfo, err := os.Stat(f.path(name))
	if err != nil {
		return 0, err
	}
	return finfo.Mode(), nil
}
func (f *fsPerm) ChOwner(name string, owner string) error {
	user, err := user.Lookup(owner)
	if err != nil {
		return err
	}
	uid, _ := strconv.Atoi(user.Uid)
	return os.Chown(f.path(name), uid, -1)
}
func (f *fsPerm) ChGroup(name string, groupName string) error {
	group, err := user.Lookup(groupName)
	if err != nil {
		return err
	}
	gid, _ := strconv.Atoi(group.Gid)
	return os.Chown(f.path(name), -1, gid)
}
func (f *fsPerm) ChMode(name string, mode os.FileMode) error {
	return os.Chmod(f.path(name), mode)
}

