package httpd

import (
	"archive/tar"
	"bytes"
	"compress/bzip2"
	_ "embed"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

//go:embed icons.tar.bz2
var iconsArchive []byte
var icons map[string][]byte
var initTime time.Time

// Return the location of an icon to display for a file of a given
// mime-type within iwebd's HTTP server.
func iconForMimeType(mimeType string) string {
	if pos := strings.Index(mimeType, ";"); pos >= 0 {
		mimeType = mimeType[:pos]
	}
	fileCandidate := "icons/" + strings.Replace(mimeType, "/", "-", -1) + ".png"
	if getNamedIcon(fileCandidate) != nil {
		return "/.well-known/" + fileCandidate
	}
	return ""
}

// Return the location of an icon to display for a file,
// given the fs.FileInfo object for the file and a function to retrieve
// the first few bytes of the file.
func getIconURI(file fs.FileInfo, headerGetter func() []byte) string {
	if file.IsDir() {
		return "/.well-known/icons/folder.png"
	}

	if file.Size() == 0 {
		return "/.well-known/icons/empty.png"
	}

	ext := filepath.Ext(file.Name())
	mimeType := mime.TypeByExtension(ext)

	if pos := strings.Index(mimeType, ";"); pos >= 0 {
		mimeType = mimeType[:pos]
	}

	iconPath := iconForMimeType(mimeType)
	if iconPath != "" {
		return iconPath
	}

	mimeType = http.DetectContentType(headerGetter())
	iconPath = iconForMimeType(mimeType)
	if iconPath != "" {
		return iconPath
	}

	return "/.well-known/icons/unknown.png"
}

// Load a given icon from the embedded file icon archive.
func getNamedIcon(icon string) (io.ReadSeeker) {
	if icons == nil {
		compressedArchiveReader := bytes.NewReader(iconsArchive)
		archiveReader := bzip2.NewReader(compressedArchiveReader)
		tarReader := tar.NewReader(archiveReader)

		icons = make(map[string][]byte)

		for {
			header, err := tarReader.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				panic(err)
			}
			if header.Size == 0 {
				continue
			}
			data := make([]byte, header.Size)
			r, err := tarReader.Read(data)
			if int64(r) != header.Size || (err != io.EOF && err != nil) {
				panic(err)
			}
			icons[header.Name] = data
		}

		initTime = time.Now()
	}

	val, ok := icons[icon]
	if ok {
		return bytes.NewReader(val)
	} else {
		return nil
	}
}

// Serve one of iwebd's file type images
//
// This expects paths /.well-known/icons/foo.png, and will serve the specified file from
// the icons archive.
func ServeIcon(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Not Allowed", 400)
		return
	}
	reader := getNamedIcon(r.URL.Path[len("/.well-known/"):])
	if reader == nil {
		http.NotFound(w, r)
		return
	}
	http.ServeContent(w, r, "file.png", initTime, reader)
}
