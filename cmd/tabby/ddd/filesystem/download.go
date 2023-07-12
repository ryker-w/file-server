package filesystem

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/lishimeng/file-server/internal/etc"
	"github.com/lishimeng/go-log"
	"os"
)

type downloadFileReq struct {
	FilePath string `json:"filePath,omitempty"`
	FileName string `json:"fileName,omitempty"`
}

func download(ctx iris.Context) {
	var err error
	var req downloadFileReq
	err = ctx.ReadJSON(&req)
	if err != nil {
		log.Debug(err)
	}
	var path = etc.Config.FileSystem.Root + req.FilePath
	bytes, _, err := fileToBytes(path)
	if err != nil {
		log.Debug(err)
	}
	fmt.Printf("read fileName %s", req.FileName)
	//设置文件类型
	ctx.Header("Content-Type", "application/octet-stream;charset=utf8")
	//设置文件名称
	ctx.Header("Content-Disposition", "attachment;filename="+req.FileName)
	ctx.ContentType("application/octet-stream")
	ctx.Write(bytes)
}

// 读取文件到[]byte中
func fileToBytes(filename string) ([]byte, *os.File, error) {
	// File
	file, err := os.Open(filename)
	if err != nil {
		return nil, file, err
	}
	defer file.Close()
	// FileInfo:
	stats, err := file.Stat()
	if err != nil {
		return nil, file, err
	}
	// []byte
	data := make([]byte, stats.Size())
	count, err := file.Read(data)
	if err != nil {
		return nil, file, err
	}
	fmt.Printf("read file %s len: %d \n", filename, count)
	return data, file, nil
}
