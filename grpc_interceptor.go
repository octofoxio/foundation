/*
 * Copyright (c) 2019. Octofox.io
 */

package foundation

import (
	"context"
	"fmt"
	foundationerrorv2 "github.com/octofoxio/foundation/errors/v2"
	"github.com/octofoxio/foundation/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func PanicRecoveryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				switch e := r.(type) {
				case *foundationerrorv2.Error:
					err = status.Error(e.Type, e.Error())
				case error:
					err = status.Error(codes.Unknown, err.Error())
				case string:
					err = status.Error(codes.Unknown, e)
				default:
					fmt.Println("======== Unknown error occurs =========")
					fmt.Println(r)
					panic(r)
				}
			}
		}()
		return handler(ctx, req)
	}
}

// logging a method calling information
func WithMethodCallingLoggerServerInterceptor(logger *logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		l := logger.WithServiceInfo(info.FullMethod)
		requestID := GetRequestIDFromContext(ctx)
		l = l.WithRequestID(requestID)
		l.Println("calling: " + info.FullMethod)
		if r, ok := req.(fmt.Stringer); ok {
			l.Println("Body: " + r.String())
		}
		resp, err = handler(ctx, req)
		if resp != nil {
			if r, ok := resp.(fmt.Stringer); ok {
				l.Println("Body: " + r.String())
			}
		}
		return resp, err
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
