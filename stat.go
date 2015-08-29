// Package times provides a platform-independent way to get atime, mtime, ctime and btime for files.
package times

import (
	"os"
	"time"
)

// Get returns the Timespec for the given FileInfo
func Get(fi os.FileInfo) Timespec {
	return getTimespec(fi)
}

// Stat returns the Timespec for the given filename
func Stat(name string) (Timespec, error) {
	fi, err := os.Stat(name)
	if err != nil {
		return nil, err
	}
	return getTimespec(fi), nil
}

// Timespec provides access to file times.
// ChangeTime() and BirthTime() return true if and only if they are available on the platform.
type Timespec interface {
	ModTime() time.Time
	AccessTime() time.Time
	ChangeTime() (time.Time, bool)
	BirthTime() (time.Time, bool)
}

type atime struct {
	v time.Time
}

func (a atime) AccessTime() time.Time { return a.v }

type ctime struct {
	v time.Time
}

type mtime struct {
	v time.Time
}

func (m mtime) ModTime() time.Time { return m.v }

func (c ctime) ChangeTime() (time.Time, bool) { return c.v, true }

type btime struct {
	v time.Time
}

func (b btime) BirthTime() (time.Time, bool) { return b.v, true }

type noctime struct{}

func (c noctime) ChangeTime() (t time.Time, ok bool) { return t, false }

type nobtime struct{}

func (c nobtime) BirthTime() (t time.Time, ok bool) { return t, false }
