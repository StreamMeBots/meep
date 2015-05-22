package main

import (
	"fmt"
	assetfs "github.com/elazarl/go-bindata-assetfs"

	"net/http"
)

type fileSystem struct {
	http.FileSystem
}

func (fs *fileSystem) Open(name string) (http.File, error) {
	fmt.Println("open", name)
	f, err := fs.FileSystem.Open(name)
	if err != nil {
		return fs.FileSystem.Open("index.html")
	}
	return f, nil
}

func (b *fileSystem) Exists(prefix string, filepath string) bool {
	return true
}

func Assets() *fileSystem {
	return &fileSystem{&assetfs.AssetFS{
		Asset:    Asset,
		AssetDir: AssetDir,
		Prefix:   "client",
	}}
}
