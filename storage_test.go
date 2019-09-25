/*
 * Copyright (c) 2019. Octofox.io
 */

package foundation

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/octofoxio/foundation/errors"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
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

	readCloser, err := local.GetObjectReader("./storage.go")
	t.Log(err)
	defer func() { _ = readCloser.Close() }()
	b, err := ioutil.ReadAll(readCloser)
	t.Log(err)
	t.Log(string(b))
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

func TestS3StorageIntegration(t *testing.T) {
	// only on credential provide
	// and require create a bucket name foundation-test, to run full test
	// * Access Key ID:     AWS_ACCESS_KEY_ID
	// * Secret Access Key: AWS_SECRET_ACCESS_KEY
	if os.Getenv("AWS_ACCESS_KEY_ID") == "" {
		t.Skip()
	}
	awsConfig := &aws.Config{
		Region: aws.String("ap-southeast-1"),
		Credentials: credentials.NewChainCredentials([]credentials.Provider{
			&credentials.EnvProvider{},
		}),
	}
	ss := NewS3FileStorage("foundation-test", awsConfig)

	err := ss.PutObject(".foundationrc", []byte("just for fun"))
	assert.NoError(t, err)

	rcFile, err := ss.GetObjectURL(".foundationrc")
	assert.NoError(t, err)
	t.Log("get object url: " + rcFile)

	t.Log("test invalid key")
	{
		rcFile, err := ss.GetObjectURL(".foundationrcc")
		assert.Error(t, err)
		assert.NotEmpty(t, rcFile)
		assert.IsType(t, &errors.Error{}, err)
		assert.EqualValues(t, err.(*errors.Error).Type(), errors.ErrorTypeNotfound)
		t.Log("error: ", err.Error())
	}

	t.Log("test get object")
	{
		output, err := ss.GetObject(".foundationrc")
		assert.NoError(t, err)
		assert.NotEmpty(t, output)
		t.Logf("body: %s", string(output))
	}
}
