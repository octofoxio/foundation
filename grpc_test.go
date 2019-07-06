/*
 * Copyright (c) 2019. Octofox.io
 */

package foundation

import (
	"context"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
	"path"
	"testing"
	"time"
)

type TestService2 struct {
	requestID string
	token     string
}

func (t *TestService2) Ping(c context.Context, input *PingInput) (*PingOutput, error) {
	t.requestID = GetRequestIDFromContext(c)
	return &PingOutput{}, nil
}

type TestService struct {
	requestID string
	token     string
}

func (t *TestService) Ping(c context.Context, input *PingInput) (*PingOutput, error) {
	conn := MakeDialOrPanic("localhost:3082")
	cc := NewTestClient(conn)
	_, _ = cc.Ping(c, input)
	t.requestID = GetRequestIDFromContext(c)
	if token := GetAccessTokenFromContext(c); token != "" {
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
	serv2 := NewGRPCServer()
	testservice := &TestService{}
	RegisterTestServer(serv, testservice)
	testservice2 := &TestService2{}
	RegisterTestServer(serv2, testservice2)

	reflection.Register(serv)
	lis, err := net.Listen("tcp", "localhost:3081")
	lis2, err := net.Listen("tcp", "localhost:3082")
	assert.Nil(t, err)
	go func() {
		_ = serv.Serve(lis)
	}()
	go func() {
		_ = serv2.Serve(lis2)
	}()
	time.Sleep(1 * time.Second)
	conn, err := MakeDial("localhost:3081")
	client := NewTestClient(conn)

	assert.Nil(t, err)
	output, err := client.Ping(context.Background(), &PingInput{
		Greeting: "Ho",
	})
	t.Run("TLS test", func(t *testing.T) {
		assert.Nil(t, err)
		assert.Equal(t, output.Greeting, "HI")
	})

	t.Run("WithSession", func(t *testing.T) {
		ctx := AppendAuthorizationToContext(context.Background(), "TESTTTT")
		output, err = client.Ping(ctx, &PingInput{
			Greeting: "Ho",
		})
		assert.Nil(t, err)
		assert.Equal(t, output.Greeting, "TESTTTT")
		assert.NotEmpty(t, testservice.requestID)
	})

	t.Run("RequestID should able to send via context", func(t *testing.T) {
		c := context.Background()
		_, _ = client.Ping(c, &PingInput{})
		assert.NotEmpty(t, testservice.requestID)
		assert.NotEmpty(t, testservice2.requestID)
		assert.EqualValues(t,
			testservice.requestID,
			testservice2.requestID,
		)
	})

	serv.Stop()
}
