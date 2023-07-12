package filesystem

import (
	"crypto/md5"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/lishimeng/app-starter"
	"github.com/lishimeng/app-starter/tool"
	"github.com/lishimeng/file-server/internal/etc"
	"github.com/lishimeng/file-server/internal/utils"
	"github.com/lishimeng/go-log"
	"github.com/pkg/errors"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

const (
	maxSize         = 1 * iris.GB
	uploadFormParam = "path"
)

type UploadResp struct {
	app.Response
}

// upload @[]path: 路径, @[]file: 文件
func upload(ctx iris.Context) {

	var err error
	var resp UploadResp
	ctx.SetMaxRequestBodySize(maxSize)

	folder := fmt.Sprintf("%d", time.Now().Unix())
	root, err := os.MkdirTemp(folder, "upload-*")
	if err != nil {
		log.Debug(errors.Wrapf(err, "create temp folder fail: %s", folder))
		resp.Code = tool.RespCodeError
		tool.ResponseJSON(ctx, resp)
		return
	}
	defer func() { // 清理缓存
		_ = os.RemoveAll(root)
	}()

	uploaded, _, err := ctx.UploadFormFiles(root, beforeSave)
	if err != nil {
		log.Debug(err)
		resp.Code = iris.StatusRequestEntityTooLarge
		tool.ResponseJSON(ctx, resp)
		return
	}

	form := ctx.Request().MultipartForm
	path := form.Value[uploadFormParam]
	if len(path) == 0 {
		log.Debug(errors.New("path nil, skip"))
		resp.Code = tool.RespCodeSuccess
		tool.ResponseJSON(ctx, resp)
		return
	}
	relPath := filepath.Join(path...)
	log.Debug("request path:%s", relPath)
	destDir := filepath.Join(etc.Config.FileSystem.Root, relPath)
	log.Debug("prepare to save files to:%s", destDir)

	err = os.MkdirAll(destDir, 0755)
	if err != nil {
		log.Debug(errors.Wrapf(err, "create dest dir fail:%s", destDir))
		resp.Code = tool.RespCodeError
		tool.ResponseJSON(ctx, resp)
		return
	}

	var files []string
	for _, u := range uploaded { // 存储文件名
		files = append(files, u.Filename)
	}

	for _, f := range files { // 保存文件:复制到指定文件夹
		src := filepath.Join(root, f)
		dest := filepath.Join(destDir, f)
		err = utils.CopyFile(src, dest)
		if err != nil {
			log.Debug(errors.Wrapf(err, "save file fail:%s", dest))
			resp.Code = tool.RespCodeError
			tool.ResponseJSON(ctx, resp)
			return
		}
	}

	resp.Code = tool.RespCodeSuccess
	tool.ResponseJSON(ctx, resp)
}

func beforeSave(_ iris.Context, file *multipart.FileHeader) bool {
	// TODO
	h := md5.New()
	h.Write([]byte(file.Filename))
	bs := h.Sum(nil)
	file.Filename = string(bs)
	log.Debug("save file: %s\n", file.Filename)
	return true
}
