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
		t.Errorf("expected atime to be recent: got %v instead of ~%v", at, et)
	}
	if at.ModTime().Before(et) {
		t.Errorf("expected mtime to be recent: got %v instead of ~%v", at, et)
	}
	if ct, ok := at.ChangeTime(); ok && ct.Before(et) {
		t.Errorf("expected ctime to be recent: got %v instead of ~%v", at, et)
	}
	if bt, ok := at.BirthTime(); ok && bt.Before(et) {
		t.Errorf("expected btime to be recent: got %v instead of ~%v", at, et)
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
		t.Errorf("expected atime to be recent: got %v instead of ~%v", at, et)
	}
	if at.ModTime().Before(et) {
		t.Errorf("expected mtime to be recent: got %v instead of ~%v", at, et)
	}
	if ct, ok := at.ChangeTime(); ok && ct.Before(et) {
		t.Errorf("expected ctime to be recent: got %v instead of ~%v", at, et)
	}
	if bt, ok := at.BirthTime(); ok && bt.Before(et) {
		t.Errorf("expected btime to be recent: got %v instead of ~%v", at, et)
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
	nc.ChangeTime()

	var nb nobtime
	nb.BirthTime()
}
