/*
 * Copyright (c) 2019. Octofox.io
 */

package foundationerrorv2

import (
	"google.golang.org/grpc/codes"
	"testing"
)

func TestError(t *testing.T) {

	err := New(codes.NotFound)
	panic(err)
}
