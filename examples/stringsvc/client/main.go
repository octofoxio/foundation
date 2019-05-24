/*
 * Copyright (c) 2019. Octofox.io
 */

package main

import (
	"context"
	"fmt"
	"github.com/octofoxio/foundation/examples/stringsvc/app"
	"github.com/octofoxio/foundation/grpc"
)

func main() {

	conn := grpc.MakeDialOrPanic("localhost:3010")
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
