/*
 * Copyright (c) 2019. Octofox.io
 */

package fs

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/octofoxio/foundation/logger"
	"io"
	"io/ioutil"
	"os"
	"path"
	"time"
)

type FileStorage interface {
	GetObject(key string) (result []byte, err error)
	PutObject(key string, data []byte) (err error)
	RemoveObject(key string) (err error)
	GetJSONObject(key string, data interface{}) (err error)
	GetObjectURL(key string) (url string, err error)
	GetObjectReader(key string) (result io.ReadCloser, err error)
	GetPreSignUploadURL(key string, size int64) (url string, err error)
}

type S3FileStorage struct {
	BucketName string
	awsConfig  *aws.Config
	log        *logger.Logger
}

func (s *S3FileStorage) GetJSONObject(key string, data interface{}) (err error) {
	file, err := s.GetObject(key)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(file, data); err != nil {
		return err
	}
	return nil
}

func (s *S3FileStorage) GetObjectURL(key string) (url string, err error) {
	s3Client, err := s.s3()
	if err != nil {
		return
	}
	request, _ := s3Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(key),
	})
	url, err = request.Presign(10 * time.Minute)
	return
}

func (s *S3FileStorage) RemoveObject(key string) (err error) {
	s3Client, err := s.s3()
	output, err := s3Client.DeleteObject(&s3.DeleteObjectInput{
		Key:    aws.String(key),
		Bucket: aws.String(s.BucketName),
	})
	if err != nil {
		return err
	}
	s.log.Printf("Remove fileInfo from storage %s", output.String())
	return nil
}

func (s *S3FileStorage) GetPreSignUploadURL(key string, size int64) (url string, err error) {
	s3Client, err := s.s3()
	if err != nil {
		return
	}

	request, _ := s3Client.PutObjectRequest(&s3.PutObjectInput{
		Bucket:        aws.String(s.BucketName),
		Key:           aws.String(key),
		ContentLength: aws.Int64(size),
	})

	url, err = request.Presign(10 * time.Minute)
	return
}

func (s *S3FileStorage) GetObjectReader(key string) (reader io.ReadCloser, err error) {
	s3Client, err := s.s3()
	if err != nil {
		return
	}
	output, err := s3Client.GetObject(&s3.GetObjectInput{
		Key:    aws.String(key),
		Bucket: aws.String(s.BucketName),
	})
	if err != nil {
		return
	}
	reader = output.Body
	return
}

func (s *S3FileStorage) s3() (s3Client *s3.S3, err error) {
	awsSession, err := session.NewSession(s.awsConfig)
	if err != nil {
		return
	}
	s3Client = s3.New(awsSession)
	return
}

func (s *S3FileStorage) PutObject(key string, data []byte) (err error) {

	s3Client, err := s.s3()
	if err != nil {
		return err
	}

	r := bytes.NewReader(data)
	putObjectOutput, err := s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(key),
		Body:   r,
	})
	if err != nil {
		return err
	}
	s.log.Printf("S3FileStorage: fileInfo upload complete, %s", *putObjectOutput.ETag)
	return

}

func (s *S3FileStorage) GetObject(key string) (result []byte, err error) {
	output, err := s.GetObjectReader(key)
	if err != nil {
		return result, err
	}
	result, err = ioutil.ReadAll(output)
	return
}

func NewS3FileStorage(bucketName string, awsConfig *aws.Config) *S3FileStorage {
	return &S3FileStorage{BucketName: bucketName, awsConfig: awsConfig, log: logger.New("S3FileStorage")}
}

type LocalFileStorage struct {
	Path string
}

func (l *LocalFileStorage) GetJSONObject(key string, data interface{}) (err error) {
	file, err := l.GetObject(key)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(file, data); err != nil {
		return err
	}
	return nil
}

// GetObjectURL require to implement path for it (local only)
// see dev_api on local storage URL support
func (l *LocalFileStorage) GetObjectURL(key string) (url string, err error) {
	if _, err := os.Stat(path.Join(l.Path, key)); os.IsNotExist(err) {
		return "", nil
	}
	return fmt.Sprintf("/assets/%s", key), nil
}

func (l *LocalFileStorage) RemoveObject(key string) (err error) {
	fmt.Println(key)
	if _, err := os.Stat(path.Join(l.Path, key)); os.IsNotExist(err) {
		return nil
	} else {
		return os.Remove(path.Join(l.Path, key))
	}
}

func NewLocalFileStorage(path string) *LocalFileStorage {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err = os.MkdirAll(path, 0755); err != nil {
			panic(err)
		}
	}
	return &LocalFileStorage{
		Path: path,
	}
}
func (l *LocalFileStorage) GetPreSignUploadURL(key string, size int64) (url string, err error) {
	panic("implement me")
}

func (l *LocalFileStorage) GetObjectReader(key string) (result io.ReadCloser, err error) {
	panic("implement me")
}

func (l *LocalFileStorage) PutObject(key string, data []byte) (err error) {
	filePath := path.Join(l.Path, key)
	{
		if _, e := os.Stat(path.Dir(filePath)); os.IsNotExist(e) {
			err = os.MkdirAll(path.Dir(filePath), os.ModePerm)
			if err != nil {
				return err
			}
		}
		f, err := os.Create(filePath)
		if err != nil {
			return err
		}
		w := bufio.NewWriter(f)
		_, err = w.Write(data)
		if err != nil {
			return err
		}
		err = w.Flush()
		if err != nil {
			return err
		}
	}

	return err
}

func (l *LocalFileStorage) GetObject(key string) (result []byte, err error) {
	result, err = ioutil.ReadFile(path.Join(l.Path, key))
	if err != nil {
		return
	}
	return
}
