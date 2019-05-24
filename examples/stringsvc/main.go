/*
 * Copyright (c) 2019. Octofox.io
 */

package main

import (
	"context"
	"encoding/json"
	"github.com/octofoxio/foundation/examples/stringsvc/app"
	"github.com/octofoxio/foundation/grpc"
	"github.com/octofoxio/foundation/http"
	"github.com/octofoxio/foundation/logger"
	"net"
	http2 "net/http"
)

func makeConcatRequestDecoder() http.RequestDecoder {
	return func(ctx context.Context, r *http2.Request) (i interface{}, e error) {
		p1 := r.FormValue("param1")
		p2 := r.FormValue("param2")
		return &app.ConcatInput{
			Origin: p1,
			Extend: p2,
		}, nil
	}
}

func makeConcatResponseEncoder() http.ResponseEncoder {
	return func(ctx context.Context, response interface{}) (i int, bytes []byte, e error) {
		var res = response.(*app.ConcatOutput)
		b, err := json.Marshal(res)
		if err != nil {
			return 500, nil, err
		}
		return 200, b, nil
	}
}
func main() {
	var log = logger.New("stringsvc").WithServiceInfo("main")
	stringsvc := app.NewStringSvc()

	grpcServer := grpc.NewGRPCServer()
	app.RegisterStringServer(grpcServer, stringsvc)
	go func() {
		lis, err := net.Listen("tcp", "0.0.0.0:3010")
		if err != nil {
			panic(err)
		}
		if err := grpcServer.Serve(lis); err != nil {
			panic(err)
		}
	}()

	httpServer := http.NewServer()
	httpServer.Get("/concat",
		http.EndpointHandler(func(ctx context.Context, request interface{}) (i interface{}, e error) {
			return stringsvc.Concat(ctx, request.(*app.ConcatInput))
		}),
		http.RequestDecoderMiddleware(makeConcatRequestDecoder()),
		http.ResponseEncoderMiddleware(makeConcatResponseEncoder()),
	)

	log.Info("HTTP Stringsvc start at :3009")
	log.Println("Try it on http://localhost:3009/concat?param1=hello&param2=world")
	log.Info("GRPC Stringsvc start at :3010")
	log.Println("Try it on ./client")
	err := http2.ListenAndServe("0.0.0.0:3009", httpServer)
	if err != nil {
		panic(err)
	}
}
