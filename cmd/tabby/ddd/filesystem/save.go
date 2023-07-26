package filesystem

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/lishimeng/app-starter"
	"github.com/lishimeng/app-starter/tool"
	"github.com/lishimeng/file-server/internal/etc"
	"github.com/lishimeng/go-log"
	"github.com/pkg/errors"
	"io"
	"os"
	"path/filepath"
)

type SaveFileReq struct {
	Path []string `json:"path,omitempty"` // 文件夹路径
	Name string   `json:"name,omitempty"` // 名称, 如果空,会使用src的名称覆盖
	Src  string   `json:"src,omitempty"`  // 原文件
	Hash bool     `json:"hash,omitempty"` // 是否使用hash名称, 会覆盖Name
}

func save(ctx iris.Context) {

	var err error
	var resp app.Response
	var req SaveFileReq

	err = ctx.ReadJSON(&req)
	if err != nil {
		log.Debug(err)
		resp.Code = tool.RespCodeError
		return
	}
	var src = filepath.Join(os.TempDir(), req.Src)
	_, err = os.Stat(src)
	if os.IsNotExist(err) {
		if err != nil {
			log.Debug(err)
			resp.Code = tool.RespCodeNotFound
			resp.Message = "not found src:" + req.Src
			return
		}
	}

	root := etc.Config.FileSystem.Root
	dir := buildDir(root, req.Path...)
	var name = req.Name
	if len(req.Name) == 0 {
		name = filepath.Base(req.Src)
	}
	var ext = filepath.Ext(name)
	if req.Hash {
		name = fmt.Sprintf("%s%s", tool.UUIDString(), ext)
	}

	go _save(src, dir, name)
}

func _save(src string, dir string, dest string) {
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		log.Debug("create dir: %s", dir)
		err = os.MkdirAll(filepath.Base(dir), 0750)
	}
	if err != nil {
		log.Debug(errors.Wrapf(err, "create dir fail: %s", dir))
		return
	}
	srcFile, err := os.Open(src)
	if err != nil {
		log.Debug(errors.Wrapf(err, "can't open src: %s", src))
		return
	}
	defer func() {
		_ = srcFile.Close()
	}()

	destFile, err := os.Create(filepath.Join(dir, dest))
	if err != nil {
		log.Debug(errors.Wrapf(err, "can't create dest: %s", dest))
		return
	}

	defer func() {
		_ = destFile.Close()
	}()
	_, err = io.Copy(srcFile, destFile)
	if err != nil {
		log.Debug(errors.Wrapf(err, "can't copy [%s -> %s]", src, dest))
		return
	}
}
