/*
 * Copyright (c) 2019. Octofox.io
 */

package http

import (
	"context"
	"github.com/octofoxio/foundation"
	"net/http"
)

type RequestDecoder func(ctx context.Context, r *http.Request) (interface{}, error)

func RequestDecoderMiddleware(decoder RequestDecoder) Middleware {
	return func(next Handler) Handler {
		return func(ctx context.Context) (i interface{}, e error) {
			req := getRequestFromContext(ctx)
			log := foundation.GetLoggerFromContext(ctx)
			input, err := decoder(ctx, req)
			if err != nil {
				log.WithError(err).Warn("Request decoding error")
				return nil, err
			}
			ctx = context.WithValue(ctx, InputBodyContextKey, input)
			return next(ctx)
		}
	}
}

type ResponseEncoder func(ctx context.Context, response interface{}) (int, []byte, error)

func ResponseEncoderMiddleware(encoder ResponseEncoder) Middleware {
	return func(next Handler) Handler {
		return func(ctx context.Context) (i interface{}, e error) {
			output, err := next(ctx)
			w := getResponseWriterFromContext(ctx)
			code, body, err := encoder(ctx, output)
			if err != nil {
				return nil, err
			}
			w.WriteHeader(code)
			_, err = w.Write(body)
			if err != nil {
				return nil, err
			}
			return output, err
		}
	}
}
