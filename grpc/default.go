/*
 * Copyright (c) 2019. Octofox.io
 */

package grpc

import (
	_ "github.com/rs/xid"
)

const (
	// Generate TLS certification for GRPC
	// ```
	// docker run -itv `pwd`:/opt/certstrap/out biscarch/certstrap init --common-name "*" (For wildcard certificate)
	// mv *.crt server_test.crt
	// mv *.key server_test.key
	// ```
	OCTOFOX_FOUNDATION_GRPC_CERT = "OCTOFOX_FOUNDATION_GRPC_CERT"
	OCTOFOX_FOUNDATION_GRPC_KEY  = "OCTOFOX_FOUNDATION_GRPC_KEY"

	GRPC_METADATA_AUTHORIZATION_KEY = "Authorization"
	GRPC_METADATA_REQUEST_ID_KEY    = "RequestID"
)
