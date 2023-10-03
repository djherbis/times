package times

import (
	"errors"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"golang.org/x/sys/unix"
)

func timeToStatx(t time.Time) unix.StatxTimestamp {
	nsec := time.Duration(t.UnixNano()) * time.Nanosecond
	nsec -= time.Duration(t.Unix()) * time.Second
	return unix.StatxTimestamp{Sec: t.Unix(), Nsec: uint32(nsec)}
}

func statxT(t time.Time, hasBtime bool) *unix.Statx_t {
	var statx unix.Statx_t

	statxt := timeToStatx(t)

	statx.Atime = statxt
	statx.Mtime = statxt
	statx.Ctime = statxt

	if hasBtime {
		statx.Mask = unix.STATX_BTIME
		statx.Btime = statxt
	}

	return &statx
}

type statxFuncTyp func(dirfd int, path string, flags int, mask int, stat *unix.Statx_t) (err error)

func unsupportedStatx(dirfd int, path string, flags int, mask int, stat *unix.Statx_t) (err error) {
	return unix.ENOSYS
}

var errBadStatx = errors.New("bad")

func badStatx(dirfd int, path string, flags int, mask int, stat *unix.Statx_t) (err error) {
	return errBadStatx
}

func fakeSupportedStatx(ts *unix.Statx_t) statxFuncTyp {
	return func(dirfd int, path string, flags int, mask int, stat *unix.Statx_t) (err error) {
		*stat = *ts
		return nil
	}
}

func setStatx(fn statxFuncTyp) func() {
	atomic.StoreInt32(&supportsStatx, 1)
	restoreStatx := statxFunc
	statxFunc = fn
	return func() { statxFunc = restoreStatx }
}

func TestStatx(t *testing.T) {
	tests := []struct {
		name    string
		statx   statxFuncTyp
		wantErr error
	}{
		{name: "unsupported", statx: unsupportedStatx},
		{name: "fake supported with btime", statx: fakeSupportedStatx(statxT(time.Now(), true))},
		{name: "fake supported without btime", statx: fakeSupportedStatx(statxT(time.Now(), false))},
		{name: "bad stat", statx: badStatx, wantErr: errBadStatx},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Run("stat", func(t *testing.T) {
				restore := setStatx(test.statx)
				defer restore()

				fileAndDirTest(t, func(name string) {
					ts, err := Stat(name)
					if err != nil {
						if err == test.wantErr {
							return
						}
						t.Fatal(err.Error())
					}
					timespecTest(ts, newInterval(time.Now(), time.Second), t)
				})
			})

			t.Run("statFile", func(t *testing.T) {
				restore := setStatx(test.statx)
				defer restore()

				fileAndDirTest(t, func(name string) {
					fi, err := os.Open(name)
					if err != nil {
						t.Fatal(err.Error())
					}
					defer fi.Close()

					ts, err := StatFile(fi)
					if err != nil {
						if err == test.wantErr {
							return
						}
						t.Fatal(err.Error())
					}
					timespecTest(ts, newInterval(time.Now(), time.Second), t)
				})
			})

			t.Run("lstat", func(t *testing.T) {
				restore := setStatx(test.statx)
				defer restore()

				fileAndDirTest(t, func(name string) {
					ts, err := Lstat(name)
					if err != nil {
						if err == test.wantErr {
							return
						}
						t.Fatal(err.Error())
					}
					timespecTest(ts, newInterval(time.Now(), time.Second), t)
				})
			})
		})
	}
}
