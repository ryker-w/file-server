package filesystem

import (
	"github.com/kataras/iris/v12"
	"github.com/lishimeng/app-starter"
	"github.com/lishimeng/app-starter/tool"
	"github.com/lishimeng/file-server/internal/etc"
	"github.com/lishimeng/go-log"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type ExplorerReq struct {
	Path []string `json:"path,omitempty"`
}

type ExplorerResp struct {
	app.Response
	Folders []string   `json:"folders,omitempty"`
	Files   []FileInfo `json:"files,omitempty"`
}

type FileInfo struct {
	Name string `json:"name,omitempty"`
	Ext  string `json:"ext,omitempty"`
}

func explorer(ctx iris.Context) {
	var err error
	var req ExplorerReq
	var resp ExplorerResp

	err = ctx.ReadJSON(&req)
	if err != nil {
		log.Debug(err)
		resp.Code = tool.RespCodeError
		tool.ResponseJSON(ctx, resp)
		return
	}

	log.Debug(req.Path)

	var root = etc.Config.FileSystem.Root
	root, err = filepath.Abs(root)
	if err != nil {
		log.Debug(err)
		resp.Code = tool.RespCodeError
		tool.ResponseJSON(ctx, resp)
		return
	}
	dir := buildDir(root, req.Path...)
	isFolder := dirExists(dir)
	if !isFolder {
		// 不是文件夹, 或不存在, 都返回文件夹不存在
		resp.Code = tool.RespCodeNotFound
		tool.ResponseJSON(ctx, resp)
		return
	}

	resp, err = listDir(root)
	if err != nil {
		log.Debug(err)
		resp.Code = tool.RespCodeNotFound
		tool.ResponseJSON(ctx, resp)
		return
	}

	resp.Code = tool.RespCodeSuccess
	tool.ResponseJSON(ctx, resp)
}

func listDir(root string) (resp ExplorerResp, err error) {
	err = walkFolder(root, func(path string, info os.FileInfo) {
		if info.IsDir() {
			resp.Folders = append(resp.Folders, info.Name())
		} else {
			ext := filepath.Ext(path)
			ext = strings.ToLower(strings.ReplaceAll(ext, ".", ""))
			fileInfo := FileInfo{
				Name: info.Name(),
				Ext:  ext,
			}
			resp.Files = append(resp.Files, fileInfo)
		}
	})
	return
}

func walkFolder(root string, fun func(path string, info os.FileInfo)) (err error) {
	err = filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if root == path {
			return nil
		}
		fun(path, info)

		if info.IsDir() {
			return filepath.SkipDir
		} else {
			return nil
		}
	})
	return
}

func buildDir(root string, path ...string) (dir string) {
	var pathList []string
	pathList = append(pathList, root)
	pathList = append(pathList, path...)
	dir = filepath.Join(pathList...)
	return
}

func dirExists(dir string) bool {
	info, err := os.Stat(dir)
	return (err == nil || os.IsExist(err)) && info.IsDir()
}
