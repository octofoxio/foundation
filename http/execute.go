/*
 * Copyright (c) 2019. Octofox.io
 */

package http

import (
	"context"
	"github.com/octofoxio/foundation"
	"net/http"
)

func execute(method string, path string, h Handler, middleware ...Middleware) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		ctx := context.Background()
		ctx = foundation.NewContext(ctx)
		ctx = context.WithValue(ctx, RequestContextKey, request)
		ctx = context.WithValue(ctx, ResponseWriterContextKey, writer)
		ctx = foundation.AppendLoggerToContext(ctx, foundation.GetLoggerFromContext(ctx).WithURL(method, path))

		var handler = h
		for _, m := range middleware {
			handler = m(handler)
		}
		_, err := handler(ctx)
		if err != nil {
			panic(err)
		}

	}
}
