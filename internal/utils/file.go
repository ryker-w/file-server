package utils

import (
	"github.com/lishimeng/go-log"
	"io"
	"os"
)

func CopyFile(src, dest string) (err error) {
	log.Debug("prepare to save file:[%s]-->[%s]", src, dest)
	r, err := os.Open(src)
	if err != nil {
		return
	}
	defer func() {
		_ = r.Close()
	}()
	w, err := os.Create(dest)
	if err != nil {
		return
	}
	defer func() {
		_ = w.Close()
	}()
	_, err = io.Copy(w, r)
	return
}
