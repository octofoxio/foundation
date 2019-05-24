/*
 * Copyright (c) 2019. Octofox.io
 */

package grpc

import (
	"context"
	"github.com/octofoxio/foundation"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
	"path"
	"testing"
	"time"
)

type TestService struct {
	requestID string
}

func (t *TestService) Ping(c context.Context, input *PingInput) (*PingOutput, error) {
	if token := foundation.GetAccessTokenFromContext(c); token != "" {
		t.requestID = foundation.GetRequestIDFromContext(c)
		return &PingOutput{
			Greeting: token,
		}, nil
	}
	return &PingOutput{
		Greeting: "HI",
	}, nil
}

func TestNewGRPCServer(t *testing.T) {
	wd, _ := os.Getwd()
	var (
		certPath = path.Join(wd, "./server_test.crt")
		certKey  = path.Join(wd, "./server_test.key")
	)
	_ = os.Setenv(OCTOFOX_FOUNDATION_GRPC_CERT, certPath)
	_ = os.Setenv(OCTOFOX_FOUNDATION_GRPC_KEY, certKey)

	serv := NewGRPCServer()
	testservice := &TestService{}
	RegisterTestServer(serv, testservice)

	reflection.Register(serv)
	lis, err := net.Listen("tcp", "localhost:3081")
	assert.Nil(t, err)
	go func() {
		_ = serv.Serve(lis)
	}()
	time.Sleep(1 * time.Second)
	conn, err := MakeDial("localhost:3081")
	client := NewTestClient(conn)

	t.Run("TLS test", func(t *testing.T) {
		assert.Nil(t, err)
		output, err := client.Ping(context.Background(), &PingInput{
			Greeting: "Ho",
		})
		assert.Nil(t, err)
		assert.Equal(t, output.Greeting, "HI")
	})

	t.Run("WithSession", func(t *testing.T) {
		ctx := AppendAuthorizationToContext(context.Background(), "TEST")
		output, err := client.Ping(ctx, &PingInput{
			Greeting: "Ho",
		})
		assert.Nil(t, err)
		assert.Equal(t, output.Greeting, "TEST")
		assert.NotEmpty(t, testservice.requestID)
	})

	serv.Stop()
}
