/*
 * Copyright (c) 2019. Octofox.io
 */

package foundation

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewLocalFileStorage(t *testing.T) {
	wd, _ := os.Getwd()
	local := NewLocalFileStorage(wd)
	u, err := local.GetObjectURL("./storage.go")
	assert.NoError(t, err)
	assert.NotEmpty(t, u)
	t.Log(u)
}

type localAWSEndpointResolver struct{}

func (l *localAWSEndpointResolver) EndpointFor(service, region string, opts ...func(*endpoints.Options)) (endpoints.ResolvedEndpoint, error) {
	var endpoint = ""
	switch service {
	case "sns":
		endpoint = "http://localhost:4575"
	case "sqs":
		endpoint = "http://localhost:4576"
	}

	return endpoints.ResolvedEndpoint{
		URL: endpoint,
	}, nil
}

// Local AWS config
// require host to run localstack locally
// please use docker run -p "4567-4584:4567-4584" -p "8080:8080" -e DOCKER_HOST=unix:///var/run/docker.sock localstack/localstack
func NewLocalAWSConfig() *aws.Config {
	var resolver = localAWSEndpointResolver{}
	return &aws.Config{
		EndpointResolver: &resolver,
		Region:           aws.String("ap-southeast-1"),
	}
}
func TestS3FileStorage_RemoveObject(t *testing.T) {
	fs := NewS3FileStorage("", NewLocalAWSConfig())
	{
		p, err := fs.getURLByPath("test-bucket", "ap-southeast-1", "billpayment/image.png")
		assert.NoError(t, err)
		assert.EqualValues(t, "https://test-bucket.s3.ap-southeast-1.amazonaws.com/billpayment/image.png", p)
	}
	{

		p, err := fs.getURLByPath("test-bucket", "ap-southeast-1", "/billpayment/image.png")
		assert.NoError(t, err)
		assert.EqualValues(t, "https://test-bucket.s3.ap-southeast-1.amazonaws.com/billpayment/image.png", p)
	}
}
