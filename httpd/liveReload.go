package httpd

import (
	"bytes"
	_ "embed"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"strings"

	"github.com/fsnotify/fsnotify"
)

//go:embed liveReload.html
var javascriptSnippet []byte

type EmbeddedFile struct {
	*bytes.Reader
	f http.File
}

type FakeInfo struct {
	fs.FileInfo
	sz int64
}

func (f FakeInfo) Size() int64 {
	return f.sz
}

func (f EmbeddedFile) Readdir(count int) ([]fs.FileInfo, error) { return nil, nil }
func (f EmbeddedFile) Stat() (fs.FileInfo, error) {
	st, err := f.f.Stat()
	if err != nil {
		return nil, err
	}
	return FakeInfo{st, st.Size() + int64(len(javascriptSnippet))}, nil
}
func (f EmbeddedFile) Close() error { return nil }

func EmbedLiveReload(f http.File) http.File {
	data, err := io.ReadAll(f)
	if err != nil {
		f.Seek(0, io.SeekStart)
		return f
	}

	pos := bytes.Index(data, []byte("</body>"))
	if pos < 0 {
		f.Seek(0, io.SeekStart)
		return f
	}

	out := append([]byte(""), data[:pos]...)
	out = append(out, javascriptSnippet...)
	out = append(out, data[pos:]...)

	return EmbeddedFile{bytes.NewReader(out), f}
}

func LiveReloadFeed(w http.ResponseWriter, r *http.Request) {
	q, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		http.Error(w, "Not Allowed", 400)
		return
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		http.Error(w, "Not Allowed", 400)
		return
	}
	defer watcher.Close()

	for _, f := range q["f"] {
		f, err = url.QueryUnescape(f)
		if err != nil {
			http.Error(w, "Not Allowed", 400)
			return
		}
		f = strings.TrimPrefix(f, "/")
		if f == "" {
			f = "index.html"
		}
		watcher.Add(f)
	}

	w.Header().Add("Content-Type", "text/event-stream")
	w.Header().Add("Cache-Control", "no-cache")
	w.(http.Flusher).Flush()

	select {
	case <- watcher.Events:
		break
	case <- r.Context().Done():
		break
	}

	w.Write([]byte("data: done\n\n"))
}
