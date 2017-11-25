package tail

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"
)

func TestTail(t *testing.T) {
	tmp, err := ioutil.TempDir(os.TempDir(), "")
	if err != nil {
		t.Fatalf("unexpected err: %s", err)
	}
	defer func() {
		os.Remove(tmp)
	}()

	fpath := path.Join(tmp, "test")
	f, err := os.OpenFile(fpath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	if err != nil {
		t.Fatalf("unexpected err: %s", err)
	}

	go func() {
		if err = tailFile(fpath, os.Stdout); err != nil {
			t.Fatalf("unexpected err: %s", err)
		}
	}()

	timeout := time.Tick(6 * time.Second)
	num := 0
	for {
		select {
		case <-time.After(2 * time.Second):
			f.WriteString(fmt.Sprintf("test: %d\n", num))
			num++
		case <-timeout:
			return
		}
	}
}

func _TestTailNewReader(t *testing.T) {
	tmp, err := ioutil.TempDir(os.TempDir(), "")
	if err != nil {
		t.Fatalf("unexpected err: %s", err)
	}
	defer func() {
		os.Remove(tmp)
	}()

	fpath := path.Join(tmp, "test")
	_, err = os.OpenFile(fpath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	if err != nil {
		t.Fatalf("unexpected err: %s", err)
	}

	fmt.Printf("tailing file %s\n", fpath)

	if err = tailFile(fpath, os.Stdout); err != nil {
		t.Fatalf("unexpected err: %s", err)
	}
}
