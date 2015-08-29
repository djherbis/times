// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// http://golang.org/src/os/stat_nacl.go

package times

import (
	"os"
	"syscall"
	"time"
)

type timespec struct {
	atime
	mtime
	ctime
	nobtime
}

func timespecToTime(sec, nsec int64) time.Time {
	return time.Unix(sec, nsec)
}

func getTimespec(fi os.FileInfo) (t timespec) {
	stat := fi.Sys().(*syscall.Stat_t)
	t.atime.v = timespecToTime(stat.Atime, stat.AtimeNsec)
	t.mtime.v = timespecToTime(stat.Mtime, stat.MtimeNsec)
	t.ctime.v = timespecToTime(stat.Ctime, stat.CtimeNsec)
	return t
}
