/*
 * Copyright (c) 2019. Octofox.io
 */

package foundation

import (
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"time"
)

func NewGRPCServer(interceptors ...grpc.UnaryServerInterceptor) *grpc.Server {
	interceptors = append(interceptors) // panic interceptor must be implemented outside foundation

	interceptors = append([]grpc.UnaryServerInterceptor{}, interceptors...)
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
	if certPath := EnvString(OCTOFOX_FOUNDATION_GRPC_CERT, ""); certPath != "" {
		keyPath := EnvStringOrPanic(OCTOFOX_FOUNDATION_GRPC_KEY)
		creds, err := credentials.NewServerTLSFromFile(certPath, keyPath)
		if err != nil {
			panic(fmt.Errorf("could not load TLS keys: %s", err))
		}
		grpcServerOptions = append(grpcServerOptions, grpc.Creds(creds))
	}

	return grpc.NewServer(grpcServerOptions...)
}
