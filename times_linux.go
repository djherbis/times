// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// http://golang.org/src/os/stat_linux.go

package times

import (
	"errors"
	"os"
	"runtime"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
)

// HasChangeTime and HasBirthTime are true if and only if
// the target OS supports them.
const (
	HasChangeTime = true
	HasBirthTime  = false
)

type timespec struct {
	atime
	mtime
	ctime
	nobtime
}

type timespecBtime struct {
	atime
	mtime
	ctime
	btime
}

// Stat returns the Timespec for the given filename.
func Stat(name string) (Timespec, error) {
	var statx unix.Statx_t

	err := unix.Statx(unix.AT_FDCWD, name, unix.AT_EMPTY_PATH|unix.AT_STATX_SYNC_AS_STAT, unix.STATX_ATIME|unix.STATX_MTIME|unix.STATX_CTIME|unix.STATX_BTIME, &statx)
	if err != nil {
		//linux 4.10 and earlier does not support Statx syscall
		if errors.Is(err, unix.ENOSYS) {
			return stat(name, os.Stat)
		}
		return nil, err
	}

	return extractTimes(&statx), nil
}

// Lstat returns the Timespec for the given filename, and does not follow Symlinks.
func Lstat(name string) (Timespec, error) {
	var statX unix.Statx_t

	err := unix.Statx(unix.AT_FDCWD, name, unix.AT_EMPTY_PATH|unix.AT_STATX_SYNC_AS_STAT|unix.AT_SYMLINK_NOFOLLOW, unix.STATX_ATIME|unix.STATX_MTIME|unix.STATX_CTIME|unix.STATX_BTIME, &statX)
	if err != nil {
		//linux 4.10 and earlier does not support Statx syscall
		if errors.Is(err, unix.ENOSYS) {
			return stat(name, os.Lstat)
		}
		return nil, err
	}

	return extractTimes(&statX), nil
}

// StatFile returns the Timespec for the given *os.File.
func StatFile(file *os.File) (Timespec, error) {
	var statx unix.Statx_t

	err := unix.Statx(int(file.Fd()), "", unix.AT_EMPTY_PATH|unix.AT_STATX_SYNC_AS_STAT, unix.STATX_ATIME|unix.STATX_MTIME|unix.STATX_CTIME|unix.STATX_BTIME, &statx)
	if err != nil {
		//linux 4.10 and earlier does not support Statx syscall
		if errors.Is(err, unix.ENOSYS) {
			fi, err := file.Stat()
			if err != nil {
				return nil, err
			}
			return getTimespec(fi), nil
		}
		return nil, err
	}

	runtime.KeepAlive(file)
	return extractTimes(&statx), nil
}

func statxTimestampToTime(ts unix.StatxTimestamp) time.Time {
	return time.Unix(ts.Sec, int64(ts.Nsec))
}

func extractTimes(statx *unix.Statx_t) Timespec {
	if statx.Mask&unix.STATX_BTIME == unix.STATX_BTIME {
		var t timespecBtime
		t.atime.v = statxTimestampToTime(statx.Atime)
		t.mtime.v = statxTimestampToTime(statx.Mtime)
		t.ctime.v = statxTimestampToTime(statx.Ctime)
		t.btime.v = statxTimestampToTime(statx.Btime)
		return t
	}

	var t timespec
	t.atime.v = statxTimestampToTime(statx.Atime)
	t.mtime.v = statxTimestampToTime(statx.Mtime)
	t.ctime.v = statxTimestampToTime(statx.Ctime)
	return t
}

func timespecToTime(ts syscall.Timespec) time.Time {
	return time.Unix(int64(ts.Sec), int64(ts.Nsec))
}

func getTimespec(fi os.FileInfo) (t timespec) {
	stat := fi.Sys().(*syscall.Stat_t)
	t.atime.v = timespecToTime(stat.Atim)
	t.mtime.v = timespecToTime(stat.Mtim)
	t.ctime.v = timespecToTime(stat.Ctim)
	return t
}
