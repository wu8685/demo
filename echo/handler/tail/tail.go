package tail

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/wu8685/demo/echo/server"
	"github.com/wu8685/demo/echo/tools"
)

var (
	eol = []byte{'\n'}
)

func init() {
	server.Register(tail, "/tail", "GET")
}

func tail(w http.ResponseWriter, r *http.Request) {
	tmp, err := ioutil.TempDir(os.TempDir(), "")
	if err != nil {
		tools.WriteError(w, err)
	}
	defer func() {
		os.Remove(tmp)
	}()

	fpath := path.Join(tmp, "test")
	f, err := os.OpenFile(fpath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	if err != nil {
		tools.WriteError(w, err)
	}
	defer f.Close()

	go func() {
		timeout := time.Tick(60 * time.Second)
		for {
			select {
			case <-timeout:
				f.WriteString(fmt.Sprintf("timeout & stop: %s\n", time.Now().String()))
				return
			default:
				f.WriteString(fmt.Sprintf("test: %s\n", time.Now().String()))
				<-time.After(5 * time.Second)
			}
		}
	}()

	log.Printf("start tailing file: %s", fpath)

	writer := Wrap(w)
	if err := tailFile(fpath, writer); err != nil {
		tools.WriteError(w, err)
	}
}

func tailFile(fpath string, w io.Writer) error {
	f, err := os.Open(fpath)
	if err != nil {
		return fmt.Errorf("fail to open file %s: %s", fpath, err)
	}
	defer f.Close()

	var watcher *fsnotify.Watcher
	reader := bufio.NewReader(f)
	for {
		line, err := reader.ReadBytes(eol[0])
		if err != nil {
			if err != io.EOF {
				return fmt.Errorf("fail to read file %s: %s", fpath, err)
			}

			if _, err := f.Seek(-int64(len(line)), os.SEEK_CUR); err != nil {
				return fmt.Errorf("fail to seek in file %s: %s", fpath, err)
			}
			if watcher == nil {
				if watcher, err = fsnotify.NewWatcher(); err != nil {
					return fmt.Errorf("fail to new fs notify watcher: %s", err)
				}
				defer watcher.Close()

				if err = watcher.Add(fpath); err != nil {
					return fmt.Errorf("fail to add file to fs notify watcher: %s", err)
				}
			}

			waitFileUpdate(watcher, fpath)
			continue
		}

		if _, err := w.Write(line); err != nil {
			return fmt.Errorf("fail to write file content to response stream: %s", err)
		}
	}
}

func waitFileUpdate(watcher *fsnotify.Watcher, fpath string) (err error) {
	retry := 5
	for {
		select {
		case e := <-watcher.Events:
			switch e.Op {
			case fsnotify.Write:
				return nil
			default:
				return fmt.Errorf("watching unexpected operation %s on file %s", e.Op, fpath)
			}
		case err = <-watcher.Errors:
			fmt.Printf("err when watching file %s, retry remind %d: %s", fpath, retry, err)
			if retry == 0 {
				return err
			}
			retry--
			continue
		case <-time.After(30 * time.Second):
			if _, err := os.Stat(fpath); os.IsNotExist(err) {
				return fmt.Errorf("file %s not exist", fpath)
			}
		}
	}
}

func Wrap(w http.ResponseWriter) *WriteFlusher {
	wf := &WriteFlusher{
		writer: w,
	}

	if f, ok := w.(http.Flusher); ok {
		wf.flusher = f
	}

	return wf
}

type WriteFlusher struct {
	writer  io.Writer
	flusher http.Flusher
}

func (wf *WriteFlusher) Write(bs []byte) (n int, err error) {
	n, err = wf.writer.Write(bs)
	if err != nil {
		return
	}

	if wf.flusher != nil {
		wf.flusher.Flush()
	}
	return
}
