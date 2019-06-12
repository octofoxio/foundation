/*
 * Copyright (c) 2019. Octofox.io
 */

package fs

import (
	"bytes"
	"fmt"
	"github.com/rakyll/statik/fs"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"time"
)

type file struct {
	*bytes.Reader
	mif fileInfo
}

func (mf *file) Close() error { return nil } // Noop, nothing to do

func (mf *file) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil // We are not a directory but a single file
}

func (mf *file) Stat() (os.FileInfo, error) {
	return mf.mif, nil
}

type fileInfo struct {
	name string
	data []byte
}

func (mif fileInfo) Name() string       { return mif.name }
func (mif fileInfo) Size() int64        { return int64(len(mif.data)) }
func (mif fileInfo) Mode() os.FileMode  { return 0444 }        // Read for all
func (mif fileInfo) ModTime() time.Time { return time.Time{} } // Return anything
func (mif fileInfo) IsDir() bool        { return false }
func (mif fileInfo) Sys() interface{}   { return nil }

// FileSystem
// This utility help you to access local fileInfo by 2 ways
// - Local fileInfo storage
// - Bundled fileInfo system (using statik, please see rakyll/statik documentation)
type FileSystem interface {
	GetObject(key string) ([]byte, error)
	Open(name string) (http.File, error)
}

type StatikFileSystem struct {
	http.FileSystem
}

func (f *StatikFileSystem) GetObject(k string) ([]byte, error) {
	r, err := f.Open(path.Join("/", k))
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

func (f *LocalFileSystem) Open(name string) (http.File, error) {
	if name == "/" {
		name = "index.html"
	}
	b, err := f.GetObject(name)
	if err != nil {
		return nil, err
	}
	return &file{
		Reader: bytes.NewReader(b),
		mif: fileInfo{
			name: name,
			data: b,
		},
	}, nil
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
// Create FileSystem class that can access and return fileInfo
// from system fileInfo
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
			FileSystem: s,
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
