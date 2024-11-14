package httpd

import (
	"bytes"
	"io/fs"
	"net/http"
	"os"
	"path"
	"sort"
	"time"

	"facette.io/natsort"

	"go.pberndt.com/iwebd/util"
)

type fancyDirectoryIndex struct {
	parentFS http.FileSystem
	readOnly bool
}

type fileReader struct {
	bytes.Reader
}

func (b fileReader) Close() error {
	return nil
}

func (b fileReader) Readdir(int) ([]fs.FileInfo, error) {
	return nil, nil
}

type fakeStat int64

func (f fakeStat) Name() string       { return "index.html" }
func (f fakeStat) Size() int64        { return int64(f) }
func (f fakeStat) Mode() fs.FileMode  { return fs.ModeIrregular }
func (f fakeStat) ModTime() time.Time { return time.Unix(0, 0) }
func (f fakeStat) IsDir() bool        { return false }
func (f fakeStat) Sys() interface{}   { return nil }

func (b fileReader) Stat() (fs.FileInfo, error) {
	return fakeStat(b.Len()), nil
}

func (f *fancyDirectoryIndex) Open(name string) (http.File, error) {
	dirName, baseName := path.Split(name)
	file, err := f.parentFS.Open(name)

	if err != nil && baseName == "index.html" {
		dir, err := f.parentFS.Open(dirName)
		if err != nil {
			return nil, err
		}

		files := make([]os.FileInfo, 0, 10)
		if osFile, ok := dir.(*os.File); ok {
			// Prefer the os.File interface to deal with directory entries
			// that we are not allowed to stat(2). ReadDir will error out
			// entirely at the first entry it can't handle, without being
			// able to recover.
			dirEntries, err := osFile.ReadDir(-1)
			if err != nil {
				return nil, err
			}
			for _, e := range dirEntries {
				info, err := e.Info()
				if err == nil {
					files = append(files, info)
				} else {
					util.Log.Info("Can not read %s: %v", e.Name(), err)
				}
			}
		} else {
			// Fall back to Readdir
			files, err = dir.Readdir(-1)
			if err != nil {
				return nil, err
			}
		}

		sort.Slice(files, func(ii, ji int) bool {
			i := files[ii].(fs.FileInfo)
			j := files[ji].(fs.FileInfo)
			if i.IsDir() && !j.IsDir() {
				return true
			}
			if j.IsDir() && !i.IsDir() {
				return false
			}
			return natsort.Compare(i.Name(), j.Name())
		})
		var fileList []fileContentsFragment
		for _, file := range files {
			if file.Name()[0] == byte('.') {
				continue
			}
			headerGetter := func() []byte {
				var header [512]byte
				fp, err := f.parentFS.Open(dirName + "/" + file.Name())
				if err != nil {
					return []byte{}
				}
				fp.Read(header[:])
				fp.Close()
				return header[:]
			}
			fileList = append(fileList, fileContentsFragment{
				Name:       file.Name(),
				Link:       file.Name(),
				Icon:       getIconURI(file, headerGetter),
				Annotation: util.SizeToHumanReadableSize(file.Size()),
			})
		}

		index := renderDirectoryIndexPage(dirName, fileList, f.readOnly)
		file = &fileReader{*bytes.NewReader([]byte(index))}
		return file, nil
	}


	return file, err
}
