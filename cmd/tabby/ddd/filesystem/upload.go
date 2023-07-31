package filesystem

import (
	"github.com/kataras/iris/v12"
	"github.com/lishimeng/app-starter"
	"github.com/lishimeng/app-starter/tool"
	"github.com/lishimeng/file-server/internal/etc"
	"github.com/lishimeng/file-server/internal/utils"
	"github.com/lishimeng/go-log"
	"github.com/pkg/errors"
	"net/url"
	"os"
	"path/filepath"
)

const (
	maxSize         = 1 * iris.GB
	uploadFormParam = "path"
)

type UploadResp struct {
	app.Response
	Items []UploadItemResp `json:"items,omitempty"`
}

type UploadItemResp struct {
	Path string `json:"path,omitempty"`
	Web  string `json:"web,omitempty"`
	File string `json:"file,omitempty"`
}

// upload @[]path: 路径, @[]file: 文件
func upload(ctx iris.Context) {

	var err error
	var resp UploadResp
	ctx.SetMaxRequestBodySize(maxSize)

	uploadRoot, err := os.MkdirTemp("", "upload-*")
	if err != nil {
		log.Debug(errors.Wrapf(err, "create temp folder fail"))
		resp.Code = tool.RespCodeError
		tool.ResponseJSON(ctx, resp)
		return
	}
	defer func() { // 清理缓存
		_ = os.RemoveAll(uploadRoot)
	}()

	uploaded, _, err := ctx.UploadFormFiles(uploadRoot)
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
		var destName string

		destName, err = copyFile(f, uploadRoot, destDir)
		if err != nil {
			log.Debug(errors.Wrapf(err, "save file fail:%s[%s]", destDir, destName))
			resp.Code = tool.RespCodeError
			tool.ResponseJSON(ctx, resp)
			return
		}
		responsePath := filepath.Join(relPath, destName)
		responsePath = filepath.ToSlash(responsePath)
		web, err := url.JoinPath(etc.Config.FileSystem.Domain, responsePath)
		if err != nil {
			log.Debug(errors.Wrapf(err, "save file fail:%s[%s]", destDir, destName))
			resp.Code = tool.RespCodeError
			tool.ResponseJSON(ctx, resp)
			return
		}
		resp.Items = append(resp.Items, UploadItemResp{
			Path: responsePath,
			Web:  web,
			File: f,
		})
	}

	resp.Code = tool.RespCodeSuccess
	tool.ResponseJSON(ctx, resp)
}

func copyFile(f, srcDir, destDir string) (fileName string, err error) {
	/// 名字处理
	ext := filepath.Ext(f)
	src := filepath.Join(srcDir, f)
	fileName, err = utils.FileDigest(src)
	if err != nil {
		return
	}
	fileName = fileName + ext
	///

	dest := filepath.Join(destDir, fileName)

	if err != nil {
		return
	}
	err = utils.CopyFile(src, dest)
	if err != nil {
		return
	}
	return
}
