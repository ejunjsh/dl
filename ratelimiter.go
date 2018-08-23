package main

import "time"

type ratelimiter struct {
	readNum  int64
	pasttime time.Time
	lim      int64
}

func (r *ratelimiter) wait(readNum int64) {
	if int(time.Now().UnixNano())-int(r.pasttime.UnixNano()) <= int(time.Second) {
		d := readNum - r.readNum
		if d >= r.lim {
			x := time.Second.Nanoseconds() - (time.Now().UnixNano() - r.pasttime.UnixNano())
			time.Sleep(time.Duration(x))
			r.readNum = readNum
			r.pasttime = time.Now()
		}
	} else {
		r.readNum = readNum
		r.pasttime = time.Now()
	}
}
