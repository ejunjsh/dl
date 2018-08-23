package main

import (
	"errors"
	"fmt"
	"gopkg.in/mattn/go-runewidth.v0"
	"mime"
	"net/http"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
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
		result = fmt.Sprintf("%.2fTB", float64(i)/tib)
	case i >= gib:
		result = fmt.Sprintf("%.2fGB", float64(i)/gib)
	case i >= mib:
		result = fmt.Sprintf("%.2fMB", float64(i)/mib)
	case i >= kib:
		result = fmt.Sprintf("%.2fKB", float64(i)/kib)
	default:
		result = fmt.Sprintf("%dB", i)
	}
	return
}

func formatTime(i int64) string {
	if i < 60 {
		return fmt.Sprintf("%2ds", i)
	} else if i < 3600 {
		s := i % 60
		m := i / 60
		if s == 0 {
			return fmt.Sprintf("%2dm", m)
		} else {
			return fmt.Sprintf("%2dm ", m) + formatTime(s)
		}

	} else {
		s := i % 3600
		h := i / 3600
		if s == 0 {
			return fmt.Sprintf("%2dh", h)
		} else {
			return fmt.Sprintf("%2dh ", h) + formatTime(s)
		}
	}
}

var errNoFilename = errors.New("no filename could be determined")

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

var ctrlFinder = regexp.MustCompile("\x1b\x5b[0-9]+\x6d")

func cellCount(s string) int {
	n := runewidth.StringWidth(s)
	for _, sm := range ctrlFinder.FindAllString(s, -1) {
		n -= runewidth.StringWidth(sm)
	}
	return n
}

func getLimitFromUrl(url string) (int64, string) {
	s := strings.Split(url, ":")
	if len(s) >= 2 {
		i, err := strconv.ParseInt(s[0], 0, 0)
		if err != nil {
			return -1, url
		} else {
			return i, strings.Join(s[1:], ":")
		}
	}
	return -1, url
}
