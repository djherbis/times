// +build !windows

package times

import "os"

func Stat(name string) (Timespec, error) {
	fi, err := os.Stat(name)
	if err != nil {
		return nil, err
	}
	return getTimespec(fi), nil
}
