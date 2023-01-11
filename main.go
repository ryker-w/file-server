package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/basicauth"
)

const (
	maxSize   = 1 * iris.GB
	uploadDir = "./uploads"
)

var users = map[string]string{
	"xdxxd": "dxddx",
	"ababb": "babba",
}

func init() {
	exist, err := PathExists(uploadDir)
	if !exist || err != nil {
		panic("uploads folder is not exist")
	}
	if !IsDir(uploadDir) {
		panic("uploads is not folder")
	}
}

// PathExists 判断所给路径文件/文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	//isnotexist来判断，是不是不存在的错误
	if os.IsNotExist(err) { //如果返回的错误类型使用os.isNotExist()判断为true，说明文件或者文件夹不存在
		return false, nil
	}
	return false, err //如果有错误了，但是不是不存在的错误，所以把这个错误原封不动的返回
}

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {

		return false
	}
	return s.IsDir()
}

func main() {
	app := iris.New()

	view := iris.HTML("./views", ".html")
	view.AddFunc("formatBytes", func(b int64) string {
		const unit = 1000
		if b < unit {
			return fmt.Sprintf("%d B", b)
		}
		div, exp := int64(unit), 0
		for n := b / unit; n >= unit; n /= unit {
			div *= unit
			exp++
		}
		return fmt.Sprintf("%.1f %cB",
			float64(b)/float64(div), "kMGTPE"[exp])
	})
	app.RegisterView(view)

	// Serve assets (e.g. javascript, css).
	// app.HandleDir("/public", iris.Dir("./public"))

	app.Get("/", index)

	app.Get("/upload", uploadView)
	app.Post("/upload", upload)

	filesRouter := app.Party("/files")
	{
		filesRouter.HandleDir("/", iris.Dir(uploadDir), iris.DirOptions{
			Compress: true,
			ShowList: true,

			// Optionally, force-send files to the client inside of showing to the browser.
			Attachments: iris.Attachments{
				Enable: true,
				// Optionally, control data sent per second:
				Limit: 50.0 * iris.KB,
				Burst: 1024 * iris.KB,
				// Change the destination name through:
				// NameFunc: func(systemName string) string {...}
			},

			DirList: iris.DirListRich(iris.DirListRichOptions{
				// Optionally, use a custom template for listing:
				// Tmpl: dirListRichTemplate,
				TmplName: "dirlist.html",
			}),
		})

		auth := basicauth.Default(users)

		filesRouter.Delete("/{file:path}", auth, deleteFile)
	}

	_ = app.Listen(":80")
}

func index(ctx iris.Context) {
	//ctx.Redirect("/upload")
	_ = ctx.View("index.html")
}

func uploadView(ctx iris.Context) {
	now := time.Now().Unix()
	h := md5.New()
	_, _ = io.WriteString(h, strconv.FormatInt(now, 10))
	token := fmt.Sprintf("%x", h.Sum(nil))

	if err := ctx.View("upload.html", token); err != nil {
		_, _ = ctx.HTML("<h3>%s</h3>", err.Error())
		return
	}
}

func upload(ctx iris.Context) {
	ctx.SetMaxRequestBodySize(maxSize)

	uploaded, _, err := ctx.UploadFormFiles(uploadDir, beforeSave)
	if err != nil {
		log.Println(err)
		ctx.StopWithError(iris.StatusRequestEntityTooLarge, err)
		return
	}

	for _, u := range uploaded {
		fileName := u.Filename
		log.Printf("saved [%s] \n", fileName)

		if strings.HasSuffix(fileName, ".tar.gz") || strings.HasSuffix(fileName, ".tar") {
			gzProcess(fileName)
		} else if strings.HasSuffix(fileName, ".zip") {
			_ = unzipFile(nil, "")
		}

	}

	var resp = make(map[string]interface{})
	resp["code"] = 200
	_ = ctx.JSON(resp)
}

func gzProcess(fileName string) {
	log.Printf("process gz: %s\n", fileName)
	f := getFilePath(fileName)

	defer func() {
		e := deletePath(f) // 删除原文件(压缩文件)
		if e != nil {
			log.Printf("delete gzip fail %s\n", fileName)
			log.Print(e)
		}
	}()
	handler := NewTgzPacker()
	err := handler.UnPack(f, uploadDir)
	if err != nil {
		log.Printf("tgz process failed: %s\n", fileName)
	}
}

func beforeSave(_ iris.Context, file *multipart.FileHeader) bool {
	//ip := ctx.RemoteAddr()
	//ip = strings.ReplaceAll(ip, ".", "_")
	//ip = strings.ReplaceAll(ip, ":", "_")
	//
	//file.Filename = ip + "-" + file.Filename
	log.Printf("save file: %s\n", file.Filename)
	return true
}

func getFilePath(fileName string) string {
	return path.Join(uploadDir, fileName)
}

func deletePath(filePath string) error {
	err := os.RemoveAll(filePath)
	return err
}

func deleteFile(ctx iris.Context) {
	// It does not contain the system path,
	// as we are not exposing it to the user.
	fileName := ctx.Params().Get("file")

	filePath := getFilePath(fileName)

	log.Printf("delete file: [%s]%s\n", fileName, filePath)

	if err := os.RemoveAll(filePath); err != nil {
		ctx.StopWithError(iris.StatusInternalServerError, err)
		return
	}

	ctx.Redirect("/files")
}

type TgzPacker struct {
}

func NewTgzPacker() *TgzPacker {
	return &TgzPacker{}
}

// 判断目录是否存在，在解压的逻辑会
func (tp *TgzPacker) dirExists(dir string) bool {
	info, err := os.Stat(dir)
	return (err == nil || os.IsExist(err)) && info.IsDir()
}

// UnPack tarFileName为待解压的tar包，dstDir为解压的目标目录
func (tp *TgzPacker) UnPack(tarFileName string, dstDir string) (err error) {
	// 打开tar文件
	fr, err := os.Open(tarFileName)
	if err != nil {
		return err
	}
	defer func() {
		if err2 := fr.Close(); err2 != nil && err == nil {
			err = err2
		}
	}()
	// 使用gzip解压
	if strings.HasSuffix(tarFileName, ".gz") {
		gr, err := gzip.NewReader(fr)
		if err != nil {
			return err
		}
		defer func() {
			if err2 := gr.Close(); err2 != nil && err == nil {
				err = err2
			}
		}()
		return tp.processTar(gr, dstDir)
	} else {
		return tp.processTar(fr, dstDir)
	}
}

func (tp *TgzPacker) processTar(reader io.Reader, dstDir string) error {
	// 创建tar reader
	tarReader := tar.NewReader(reader)
	// 循环读取
	for {
		header, err := tarReader.Next()
		switch {
		// 读取结束
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		case header == nil:
			continue
		}
		// 因为指定了解压的目录，所以文件名加上路径
		targetFullPath := filepath.Join(dstDir, header.Name)
		// 根据文件类型做处理，这里只处理目录和普通文件，如果需要处理其他类型文件，添加case即可
		switch header.Typeflag {
		case tar.TypeDir:
			// 是目录，不存在则创建
			if exists := tp.dirExists(targetFullPath); !exists {
				if err = os.MkdirAll(targetFullPath, 0755); err != nil {
					return err
				}
			}
		case tar.TypeReg:
			// 是普通文件，创建并将内容写入
			file, err := os.OpenFile(targetFullPath, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			_, err = io.Copy(file, tarReader)
			// 循环内不能用defer，先关闭文件句柄
			if err2 := file.Close(); err2 != nil {
				return err2
			}
			// 这里再对文件copy的结果做判断
			if err != nil {
				return err
			}
		}
	}
}

func unzipFile(file *zip.File, dstDir string) error {
	// create the directory of file
	filePath := path.Join(dstDir, file.Name)
	if file.FileInfo().IsDir() {
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			return err
		}
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return err
	}

	// open the file
	rc, err := file.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	// create the file
	w, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer w.Close()

	// save the decompressed file content
	_, err = io.Copy(w, rc)
	return err
}
