/*
 * Copyright (c) 2019. Octofox.io
 */

package foundation

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestNewLocalFileStorage(t *testing.T) {
	wd, _ := os.Getwd()
	local := NewLocalFileStorage(wd)

	u, err := GetObjectURL("./storage.go")
	assert.NoError(t, err)
	assert.NotEmpty(t, u)
	t.Log(u)

}

type Shape interface {
	NewMe() Shape
}

type Square struct {
	Width  int64
	Height int64
}

func (d *Square) NewMe() Shape {
	return &Square{
		Height: 99,
		Width:  99,
	}
}

func bind(dest interface{}) {
	d := reflect.ValueOf(dest)
	fmt.Println(d.Elem().Kind())

}

func TestArea(t *testing.T) {
	//s := &Square{Width: 10, Height: 20}
	var d Shape
	//= &Square{}
	//= Square{Width: 1}

	//bb(d)
	bind(&d)

	time.Sleep(time.Second)
}
