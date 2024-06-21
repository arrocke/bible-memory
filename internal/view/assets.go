package view

import (
	"embed"

	"github.com/benbjohnson/hashfs"
)

//go:embed assets/*
var AssetsFS embed.FS
var HashedFS = hashfs.NewFS(AssetsFS)

var ShouldHashFiles = false

func HashFile(name string) string {
    if ShouldHashFiles {
        return "/" + HashedFS.HashName(name)
    } else {
        return "/" + name
    }
}
