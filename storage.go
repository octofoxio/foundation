/*
 * Copyright (c) 2019. Octofox.io
 */

package foundation

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/octofoxio/foundation/errors"
	"github.com/octofoxio/foundation/logger"
	"io"
	"io/ioutil"
	url2 "net/url"
	"os"
	"path"
	"time"
)

type FileStorage interface {
	GetObject(key string) (result []byte, err error)
	PutObject(key string, data []byte) (err error)
	PutObjectFromReadSeeker(key string, reader io.ReadSeeker) (err error)
	PutPublicObject(key string, data []byte) (err error)
	PutPublicObjectFromReadSeeker(key string, reader io.ReadSeeker) (err error)
	RemoveObject(key string) (err error)
	GetJSONObject(key string, data interface{}) (err error)
	GetObjectURL(key string) (url string, err error)
	GetObjectPreSignURL(key string) (url string, err error)
	GetObjectReader(key string) (result io.ReadCloser, err error)
	GetPreSignUploadURL(key string, size int64) (url string, err error)
}

type S3FileStorage struct {
	BucketName string
	awsConfig  *aws.Config
	log        *logger.Logger
}

func (s *S3FileStorage) s3() (s3Client *s3.S3, err error) {
	awsSession, err := session.NewSession(s.awsConfig)
	if err != nil {
		return
	}
	s3Client = s3.New(awsSession)
	return
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

func (s *S3FileStorage) getURLByPath(bucketName, region, key string) (string, error) {
	urlParse, err := url2.Parse(fmt.Sprintf("https://%s.s3.%s.amazonaws.com/", bucketName, region))
	if err != nil {
		return "", err
	}
	urlParse.Path = path.Join(key)
	return urlParse.String(), nil
}
func (s *S3FileStorage) GetObjectURL(key string) (objectURL string, err error) {
	s3Client, err := s.s3()
	if err != nil {
		return
	}

	// create new object URL
	objectURL, err = s.getURLByPath(s.BucketName, *s3Client.Config.Region, key)
	if err != nil {
		return objectURL, err
	}

	// validate if object exists
	_, err = s3Client.GetObjectAcl(&s3.GetObjectAclInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchKey:
				return objectURL, errors.New(errors.ErrorTypeNotfound, fmt.Sprintf("file %s not found", key))
			default:
				return objectURL, err
			}
		}
		return objectURL, err
	}
	return objectURL, err
}

func (s *S3FileStorage) GetObjectPreSignURL(key string) (url string, err error) {
	s3Client, err := s.s3()
	if err != nil {
		return
	}
	request, _ := s3Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(key),
	})
	url, err = request.Presign(7 * 24 * time.Hour)
	return
}

func (s *S3FileStorage) RemoveObject(key string) (err error) {
	s3Client, err := s.s3()
	if err != nil {
		return err
	}
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

	url, err = request.Presign(15 * time.Minute)
	return
}

// GetObjectReader caller should manually close io reader
// for prevent memory leaking
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

func (s *S3FileStorage) PutObject(key string, data []byte) (err error) {
	r := bytes.NewReader(data)
	return s.PutObjectFromReadSeeker(key, r)
}

func (s *S3FileStorage) PutObjectFromReadSeeker(key string, reader io.ReadSeeker) (err error) {
	s3Client, err := s.s3()
	if err != nil {
		return err
	}

	putObjectOutput, err := s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(key),
		Body:   reader,
	})
	if err != nil {
		return err
	}
	s.log.Printf("S3FileStorage: fileInfo upload complete, %s", *putObjectOutput.ETag)
	return
}

func (s *S3FileStorage) PutPublicObject(key string, data []byte) (err error) {
	r := bytes.NewReader(data)
	return s.PutPublicObjectFromReadSeeker(key, r)
}

func (s *S3FileStorage) PutPublicObjectFromReadSeeker(key string, r io.ReadSeeker) (err error) {
	s3Client, err := s.s3()
	if err != nil {
		return err
	}

	putObjectOutput, err := s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(key),
		Body:   r,
		ACL:    aws.String("public-read"),
	})
	if err != nil {
		return err
	}
	s.log.Printf("S3FileStorage: fileInfo upload complete, %s", *putObjectOutput.ETag)
	return
}

func (s *S3FileStorage) GetObject(key string) (result []byte, err error) {
	output, err := s.GetObjectReader(key) // this method get reader from s3 API but not close
	defer func() {
		// skip close body if output == nil
		if output == nil {
			return
		}
		err = output.Close() // close response reader after read every bytes into memory
		if err != nil {
			fmt.Printf("CANNOT CLOSE OBJECT READER POSSIBLE TO HAVE SOME MEMORY LEAK %s \n", key)
			fmt.Println(err.Error())
		}
	}()
	if err != nil {
		return result, err
	}
	// get all result to byte array
	result, err = ioutil.ReadAll(output)
	return
}

func NewS3FileStorage(bucketName string, awsConfig *aws.Config) *S3FileStorage {
	return &S3FileStorage{BucketName: bucketName, awsConfig: awsConfig, log: logger.New("S3FileStorage")}
}

type LocalFileStorage struct {
	Path string
}

func (l *LocalFileStorage) PutObjectFromReadSeeker(key string, reader io.ReadSeeker) (err error) {
	r, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	return l.PutObject(key, r)
}

func (l *LocalFileStorage) PutPublicObjectFromReadSeeker(key string, reader io.ReadSeeker) (err error) {
	return l.PutObjectFromReadSeeker(key, reader)
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
	filePath := path.Join(l.Path, key)
	pathURL := url2.URL{
		Path:   filePath,
		Scheme: "file",
	}
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return pathURL.String(), err
	}
	return pathURL.String(), nil
}

func (l *LocalFileStorage) GetObjectPreSignURL(key string) (url string, err error) {
	return l.GetObjectURL(key)
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
	f, err := os.Open(key)
	if err != nil {
		return nil, err
	}
	return f, nil
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

func (l *LocalFileStorage) PutPublicObject(key string, data []byte) (err error) {
	return l.PutObject(key, data)
}

func (l *LocalFileStorage) GetObject(key string) (result []byte, err error) {
	result, err = ioutil.ReadFile(path.Join(l.Path, key))
	if err != nil {
		return
	}
	return
}
