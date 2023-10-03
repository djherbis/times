// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// http://golang.org/src/os/stat_linux.go

package times

import (
	"errors"
	"os"
	"sync/atomic"
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

var supportsStatx int32 = 1

// Stat returns the Timespec for the given filename.
func Stat(name string) (Timespec, error) {
	if atomic.LoadInt32(&supportsStatx) == 1 {
		var statx unix.Statx_t

		//https://man7.org/linux/man-pages/man2/statx.2.html
		err := unix.Statx(unix.AT_FDCWD, name, unix.AT_STATX_SYNC_AS_STAT, unix.STATX_ATIME|unix.STATX_MTIME|unix.STATX_CTIME|unix.STATX_BTIME, &statx)
		if err != nil {
			//linux 4.10 and earlier does not support Statx syscall
			if errors.Is(err, unix.ENOSYS) {
				atomic.StoreInt32(&supportsStatx, 0)
				return stat(name, os.Stat)
			}
			return nil, err
		}
		return extractTimes(&statx), nil
	}

	return stat(name, os.Stat)
}

// Lstat returns the Timespec for the given filename, and does not follow Symlinks.
func Lstat(name string) (Timespec, error) {
	if atomic.LoadInt32(&supportsStatx) == 1 {
		var statX unix.Statx_t
		//https://man7.org/linux/man-pages/man2/statx.2.html

		err := unix.Statx(unix.AT_FDCWD, name, unix.AT_STATX_SYNC_AS_STAT|unix.AT_SYMLINK_NOFOLLOW, unix.STATX_ATIME|unix.STATX_MTIME|unix.STATX_CTIME|unix.STATX_BTIME, &statX)
		if err != nil {
			//linux 4.10 and earlier does not support Statx syscall
			if errors.Is(err, unix.ENOSYS) {
				atomic.StoreInt32(&supportsStatx, 0)
				return stat(name, os.Lstat)
			}
			return nil, err
		}
		return extractTimes(&statX), nil
	}

	return stat(name, os.Lstat)
}

// StatFile returns the Timespec for the given *os.File.
func StatFile(file *os.File) (Timespec, error) {
	if atomic.LoadInt32(&supportsStatx) == 1 {
		var statx unix.Statx_t

		sc, err := file.SyscallConn()
		if err != nil {
			return nil, err
		}

		var statxErr error
		err = sc.Control(func(fd uintptr) {
			statxErr = unix.Statx(int(fd), "", unix.AT_EMPTY_PATH|unix.AT_STATX_SYNC_AS_STAT, unix.STATX_ATIME|unix.STATX_MTIME|unix.STATX_CTIME|unix.STATX_BTIME, &statx)
		})
		if err != nil {
			return nil, err
		}

		//https://man7.org/linux/man-pages/man2/statx.2.html
		if statxErr != nil {
			//linux 4.10 and earlier does not support Statx syscall
			if errors.Is(statxErr, unix.ENOSYS) {
				atomic.StoreInt32(&supportsStatx, 0)
				fi, err := file.Stat()
				if err != nil {
					return nil, err
				}
				return getTimespec(fi), nil
			}
			return nil, statxErr
		}

		return extractTimes(&statx), nil
	}

	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}
	return getTimespec(fi), nil
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
