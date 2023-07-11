package filesystem

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/lishimeng/app-starter"
	"github.com/lishimeng/app-starter/tool"
	"github.com/lishimeng/go-log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	maxSize = 1 * iris.GB
)

type UploadResp struct {
	app.Response
	FilePath string `json:"filePath,omitempty"`
}

func upload(ctx iris.Context) {
	ctx.SetMaxRequestBodySize(maxSize)

	folder := fmt.Sprintf("upload_%d", time.Now().Unix())
	root := filepath.Join(os.TempDir(), folder)

	uploaded, _, err := ctx.UploadFormFiles(root, beforeSave)
	if err != nil {
		log.Debug(err)
		ctx.StopWithError(iris.StatusRequestEntityTooLarge, err)
		return
	}

	var fileRelPath string
	for _, u := range uploaded {
		fileName := u.Filename
		log.Debug("saved [%s] \n", fileName)
		fileRelPath = strings.ReplaceAll(fileName, root, "")
	}

	var resp UploadResp
	resp.Code = tool.RespCodeSuccess
	resp.FilePath = fileRelPath
	tool.ResponseJSON(ctx, resp)
}

func beforeSave(_ iris.Context, file *multipart.FileHeader) bool {
	//ip := ctx.RemoteAddr()
	//ip = strings.ReplaceAll(ip, ".", "_")
	//ip = strings.ReplaceAll(ip, ":", "_")
	//
	//file.Filename = ip + "-" + file.Filename
	log.Debug("save file: %s\n", file.Filename)
	return true
}
