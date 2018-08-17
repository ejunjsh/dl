package main

import 	"fmt"

const (
	_KiB = 1024
	_MiB = 1048576
	_GiB = 1073741824
	_TiB = 1099511627776
)

func formatBytes(i int64) (result string) {
	switch {
	case i >= _TiB:
		result = fmt.Sprintf("%.02f TiB", float64(i)/_TiB)
	case i >= _GiB:
		result = fmt.Sprintf("%.02f GiB", float64(i)/_GiB)
	case i >= _MiB:
		result = fmt.Sprintf("%.02f MiB", float64(i)/_MiB)
	case i >= _KiB:
		result = fmt.Sprintf("%.02f KiB", float64(i)/_KiB)
	default:
		result = fmt.Sprintf("%d B", i)
	}
	return
}