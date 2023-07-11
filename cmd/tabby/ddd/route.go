package ddd

import (
	"github.com/kataras/iris/v12"
	"github.com/lishimeng/file-server/cmd/tabby/ddd/filesystem"
)

func Route(app *iris.Application) {
	root := app.Party("/api")
	filesystem.Route(root.Party("/fs"))
}
