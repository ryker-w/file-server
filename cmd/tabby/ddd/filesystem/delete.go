package filesystem

import (
	"github.com/kataras/iris/v12"
	"github.com/lishimeng/app-starter"
	"github.com/lishimeng/app-starter/tool"
	"github.com/lishimeng/file-server/internal/etc"
	"github.com/lishimeng/go-log"
	"os"
	"path/filepath"
)

type DeleteFileReq struct {
	Path []string `json:"path,omitempty"`
	Name string   `json:"name,omitempty"`
}

func deleteFile(ctx iris.Context) {
	var err error
	var req DeleteFileReq
	var resp app.Response
	err = ctx.ReadJSON(&req)
	if err != nil {
		log.Debug(err)
		resp.Code = tool.RespCodeError
		tool.ResponseJSON(ctx, resp)
		return
	}
	if len(req.Name) == 0 {
		if err != nil {
			log.Debug(err)
			resp.Code = tool.RespCodeNotFound
			tool.ResponseJSON(ctx, resp)
			return
		}
	}

	root := etc.Config.FileSystem.Root
	filePath := filepath.Join(root, req.Name)

	log.Debug("delete file: [%s]%s\n", req.Name, filePath)

	if err := os.RemoveAll(filePath); err != nil {
		ctx.StopWithError(iris.StatusInternalServerError, err)
		return
	}

	resp.Code = tool.RespCodeSuccess
	tool.ResponseJSON(ctx, resp)

}
