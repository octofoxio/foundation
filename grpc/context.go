/*
 * Copyright (c) 2019. Octofox.io
 */

package grpc

import (
	"context"
	"google.golang.org/grpc/metadata"
)

func AppendAuthorizationToContext(ctx context.Context, accessToken string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, GRPC_METADATA_AUTHORIZATION_KEY, accessToken)
}
