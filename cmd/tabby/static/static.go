package static

import (
	"embed"
)

//go:embed * assets/*/*
var Static embed.FS
