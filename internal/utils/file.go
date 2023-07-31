package utils

import (
	"crypto/sha1"
	"encoding/hex"
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

func FileDigest(path string) (digest string, err error) {
	f, err := os.Open(path)
	if err != nil {
		return digest, err
	}
	defer func() {
		_ = f.Close()
	}()

	handler := sha1.New()
	if _, err := io.Copy(handler, f); err != nil {
		return digest, err
	}

	digest = hex.EncodeToString(handler.Sum(nil))
	return
}
