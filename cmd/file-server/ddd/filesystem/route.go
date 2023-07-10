package filesystem

import "github.com/kataras/iris/v12"

func Route(root iris.Party) {
	root.Post("/explorer", explorer)
	root.Post("/download", download)
}
