// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// http://golang.org/src/os/stat_windows.go

package times

import (
	"os"
	"syscall"
	"time"
)

type timespec struct {
	atime
	mtime
	noctime
	btime
}

func getTimespec(fi os.FileInfo) (t timespec) {
	stat := fi.Sys().(*syscall.Win32FileAttributeData)
	t.atime.v = time.Unix(0, stat.LastAccessTime.Nanoseconds())
	t.mtime.v = time.Unix(0, stat.LastWriteTime.Nanoseconds())
	t.btime.v = time.Unix(0, stat.CreationTime.Nanoseconds())
	return t
}
