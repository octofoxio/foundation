/*
 * Copyright (c) 2019. Octofox.io
 */

package main

import (
	"context"
	"fmt"
	"github.com/octofoxio/foundation"
	"github.com/octofoxio/foundation/examples/stringsvc/app"
)

func main() {

	conn := foundation.MakeDialOrPanic("localhost:3010")
	stringsvcClient := app.NewStringClient(conn)
	result, err := stringsvcClient.Concat(context.Background(), &app.ConcatInput{
		Origin: "Hello",
		Extend: "World",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(result.Result)

}
