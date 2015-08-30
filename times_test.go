package times

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestStat(t *testing.T) {
	f, err := ioutil.TempFile("", "")
	et := time.Now().Add(-time.Second)
	if err != nil {
		t.Error(err.Error())
	}
	defer os.Remove(f.Name())
	defer f.Close()

	at, err := Stat(f.Name())
	if err != nil {
		t.Error(err.Error())
	}
	if at.AccessTime().Before(et) {
		t.Errorf("expected atime to be recent: got %v instead of ~%v", at.AccessTime(), et)
	}
	if at.ModTime().Before(et) {
		t.Errorf("expected mtime to be recent: got %v instead of ~%v", at.ModTime(), et)
	}
	if HasChangeTime && at.ChangeTime().Before(et) {
		t.Errorf("expected ctime to be recent: got %v instead of ~%v", at.ChangeTime(), et)
	}
	if HasBirthTime && at.BirthTime().Before(et) {
		t.Errorf("expected btime to be recent: got %v instead of ~%v", at.BirthTime(), et)
	}
}

func TestGet(t *testing.T) {
	f, err := ioutil.TempFile("", "")
	et := time.Now().Add(-time.Second)
	if err != nil {
		t.Error(err.Error())
	}
	defer os.Remove(f.Name())
	defer f.Close()

	fi, err := os.Stat(f.Name())
	if err != nil {
		t.Error(err.Error())
	}
	at := Get(fi)
	if at.AccessTime().Before(et) {
		t.Errorf("expected atime to be recent: got %v instead of ~%v", at.AccessTime(), et)
	}
	if at.ModTime().Before(et) {
		t.Errorf("expected mtime to be recent: got %v instead of ~%v", at.ModTime(), et)
	}
	if HasChangeTime && at.ChangeTime().Before(et) {
		t.Errorf("expected ctime to be recent: got %v instead of ~%v", at.ChangeTime(), et)
	}
	if HasBirthTime && at.BirthTime().Before(et) {
		t.Errorf("expected btime to be recent: got %v instead of ~%v", at.BirthTime(), et)
	}
}

func TestStatErr(t *testing.T) {
	_, err := Stat("badfile?")
	if err == nil {
		t.Error("expected an error")
	}
}

func TestCheat(t *testing.T) {
	// not all times are available for all platforms
	// this allows us to get 100% test coverage for platforms which do not have
	// ChangeTime/BirthTime
	var c ctime
	c.ChangeTime()

	var b btime
	b.BirthTime()

	var nc noctime
	func() {
		defer func() { recover() }()
		nc.ChangeTime()
	}()

	var nb nobtime
	func() {
		defer func() { recover() }()
		nb.BirthTime()
	}()
}
