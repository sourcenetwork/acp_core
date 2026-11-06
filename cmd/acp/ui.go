package main

import (
	"embed"
	"io/fs"
)

// uiFiles is the playground UI build, copied from playground2/dist by `make acp-ui`
//
//go:embed all:ui
var uiFiles embed.FS

func builtinUI() fs.FS {
	sub, err := fs.Sub(uiFiles, "ui")
	if err != nil {
		panic(err) // ui is embedded at build time: it's always there
	}
	return sub
}
