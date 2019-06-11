/*
 * Copyright (c) 2019. Octofox.io
 */

package grpc

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/octofoxio/foundation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"time"
)

func WithFoundationContext() grpc.UnaryServerInterceptor {
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
					ctx = context.WithValue(ctx, foundation.FoundationAccessTokenContextKey, token)
				}
			}

			requestIDs := md.Get(GRPC_METADATA_REQUEST_ID_KEY)
			if len(requestIDs) > 0 {
				ctx = context.WithValue(ctx, foundation.FoundationRequestIDContextKey, requestIDs[0])
			}
		}
		ctx = foundation.NewContext(ctx)
		return handler(ctx, req)
	}
}

func NewGRPCServer(interceptors ...grpc.UnaryServerInterceptor) *grpc.Server {
	interceptors = append(interceptors, grpc_recovery.UnaryServerInterceptor(
		grpc_recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			return nil
		}),
	))
	interceptors = append([]grpc.UnaryServerInterceptor{WithFoundationContext()}, interceptors...)
	var grpcServerOptions = []grpc.ServerOption{
		// To keep connection alive in-case
		// when GRPC is working
		// behind Loadbalancer
		grpc.KeepaliveParams(
			keepalive.ServerParameters{
				Time: 50 * time.Second,
			},
		),
		grpc_middleware.WithUnaryServerChain(interceptors...),
	}

	// Use TLS certification if provide to ENV
	if certPath := foundation.EnvString(OCTOFOX_FOUNDATION_GRPC_CERT, ""); certPath != "" {
		keyPath := foundation.EnvStringOrPanic(OCTOFOX_FOUNDATION_GRPC_KEY)
		creds, err := credentials.NewServerTLSFromFile(certPath, keyPath)
		if err != nil {
			panic(fmt.Errorf("could not load TLS keys: %s", err))
		}
		grpcServerOptions = append(grpcServerOptions, grpc.Creds(creds))
	}

	return grpc.NewServer(grpcServerOptions...)
}
