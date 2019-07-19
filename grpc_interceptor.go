/*
 * Copyright (c) 2019. Octofox.io
 */

package foundation

import (
	"context"
	"github.com/octofoxio/foundation/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// logging a method calling information
func WithMethodCallingLoggerServerInterceptor(logger *logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		l := logger.WithServiceInfo(info.FullMethod)
		requestID := GetRequestIDFromContext(ctx)
		l = l.WithRequestID(requestID)
		l.Infoln("method is called")
		return handler(ctx, req)
	}
}

func WithContextServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		// Try to get metadata "Authorization" from
		// request context
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			var token string
			tokens := md.Get(GRPC_METADATA_AUTHORIZATION_KEY)
			if len(tokens) > 0 {
				token = tokens[0]
				if token != "" {
					ctx = context.WithValue(ctx, FoundationAccessTokenContextKey, token)
				}
			}

			requestIDs := md.Get(GRPC_METADATA_REQUEST_ID_KEY)
			if len(requestIDs) > 0 {
				ctx = context.WithValue(ctx, FoundationRequestIDContextKey, requestIDs[0])
			}
		}
		ctx = NewContext(ctx)
		ctx = AppendRequestIDToContext(ctx, GetRequestIDFromContext(ctx))
		return handler(ctx, req)
	}
}
