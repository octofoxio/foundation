/*
 * Copyright (c) 2019. Inception Asia
 * Maintain by DigithunWorldwide ❤
 * Maintainer
 * - rungsikorn.r@digithunworldwide.com
 * - nipon.chi@digithunworldwide.com
 * - mai@digithunworldwide.com
 */

package fs

import (
	"bytes"
	"fmt"
	"github.com/rakyll/statik/fs"
	"io/ioutil"
	"net/http"
	"path"
)

// FileSystem
// This utility help you to access local file by 2 ways
// - Local file storage
// - Bundled file system (using statik, please see rakyll/statik documentation)
type FileSystem interface {
	GetObject(key string) ([]byte, error)
}

type StatikFileSystem struct {
	fs http.FileSystem
}

func (f *StatikFileSystem) GetObject(k string) ([]byte, error) {
	r, err := f.fs.Open(path.Join("/", k))
	if err != nil {
		return nil, err
	}

	var buf = new(bytes.Buffer)
	_, err = buf.ReadFrom(r)

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), err
}

type LocalFileSystem struct {
	rootDir string
}

func (f *LocalFileSystem) GetObject(key string) ([]byte, error) {
	result, err := ioutil.ReadFile(path.Join(f.rootDir, key))
	if err != nil {
		return nil, err
	}

	return result, nil
}

type StaticMode string

const (
	StaticMode_LOCAL  = "LOCAL"
	StaticMode_Statik = "STATIK"
)

// NewFileSystem
// Create FileSystem class that can access and return file
// from system file
// - rootDirectory is where the root directory of filesystem
// - mode is StaticMode_LOCAL or StaticMode_STATIK
//
// if you want to use StaticMode_STATIK please make sure you are following
// the instruction from rakll/statik documentation
func NewFileSystem(rootDir string, mode StaticMode) FileSystem {
	if mode == StaticMode_Statik {
		fmt.Println("Use fs with statik mode")
		var s, err = fs.New()
		if err != nil {
			panic(err)
		}
		return &StatikFileSystem{
			fs: s,
		}
	} else if mode == StaticMode_LOCAL {
		fmt.Println("Use fs with local mode")
		return &LocalFileSystem{
			rootDir: rootDir,
		}
	} else {
		panic("Static mode invalid")
	}
}
