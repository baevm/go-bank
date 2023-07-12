package doc

import (
	"embed"
	"io/fs"
)

//go:embed swagger/*
var SwaggerFolder embed.FS

func GetSwaggerFolder() (fs.FS, error) {
	staticFs := fs.FS(SwaggerFolder)
	sub, err := fs.Sub(staticFs, "swagger")

	return sub, err
}
