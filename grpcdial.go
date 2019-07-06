/*
 * Copyright (c) 2019. Octofox.io
 */

package foundation

import (
	"context"
	"github.com/octofoxio/foundation/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"time"
)

func MakeDialOrPanic(endpoint string) *grpc.ClientConn {
	if conn, err := MakeDial(endpoint); err != nil {
		panic(err)
	} else {
		return conn
	}
}

func WithAuthorizationDialOption(accessToken string) grpc.DialOption {
	return grpc.WithUnaryInterceptor(func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		err := invoker(ctx, method, req, reply, cc, opts...)
		ctx = metadata.AppendToOutgoingContext(ctx, GRPC_METADATA_AUTHORIZATION_KEY, accessToken)
		return err
	})
}

func MakeDial(endpoint string, dialOptions ...grpc.DialOption) (*grpc.ClientConn, error) {
	var log = logger.New("grpc").WithServiceID("foundation").WithServiceInfo("grpc")
	options := []grpc.DialOption{
		grpc.WithKeepaliveParams(
			keepalive.ClientParameters{
				Time:                50 * time.Second,
				PermitWithoutStream: true,
			}),
	}
	for _, o := range dialOptions {
		options = append(options, o)
	}

	if certPath := EnvString(OCTOFOX_FOUNDATION_GRPC_CERT, ""); certPath != "" {
		creds, err := credentials.NewClientTLSFromFile(certPath, "")
		if err != nil {
			panic(err)
		}
		options = append(options, grpc.WithTransportCredentials(creds))
	} else {
		// IF NO GRPC CERT
		// SKIP TO USE Insecure
		log.Warn("GRPC Connect dial with insecure mode")
		options = append(options, grpc.WithInsecure())
	}

	return grpc.Dial(endpoint,
		options...,
	)
}
