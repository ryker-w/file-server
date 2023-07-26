package filesystem

import "github.com/kataras/iris/v12"

func Route(root iris.Party) {
	root.Post("/explorer", explorer) // 列表
	root.Post("/download", download) // 下载

	root.Post("/upload", upload) // 上传
	root.Post("/save", save)     // 保存

	root.Post("/delete", deleteFile) // 删除
}
