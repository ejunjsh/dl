package main

import (
	"fmt"
	"net/http"
	"mime"
	"strings"
	"path/filepath"
	"path"
	"errors"
)

const (
	kib = 1024
	mib = 1048576
	gib = 1073741824
	tib = 1099511627776
)

func formatBytes(i int64) (result string) {
	switch {
	case i >= tib:
		result = fmt.Sprintf("%.02f TiB", float64(i)/tib)
	case i >= gib:
		result = fmt.Sprintf("%.02f GiB", float64(i)/gib)
	case i >= mib:
		result = fmt.Sprintf("%.02f MiB", float64(i)/mib)
	case i >= kib:
		result = fmt.Sprintf("%.02f KiB", float64(i)/kib)
	default:
		result = fmt.Sprintf("%d B", i)
	}
	return
}

var errNoFilename=errors.New("no filename could be determined")

func guessFilename(resp *http.Response) (string, error) {
	filename := resp.Request.URL.Path
	if cd := resp.Header.Get("Content-Disposition"); cd != "" {
		if _, params, err := mime.ParseMediaType(cd); err == nil {
			filename = params["filename"]
		}
	}

	if filename == "" || strings.HasSuffix(filename, "/") || strings.Contains(filename, "\x00") {
		return "", errNoFilename
	}

	filename = filepath.Base(path.Clean("/" + filename))
	if filename == "" || filename == "." || filename == "/" {
		return "", errNoFilename
	}

	return filename, nil
}