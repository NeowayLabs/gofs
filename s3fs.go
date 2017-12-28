package gofs

import (
	"bytes"
	"io"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3FS struct {
	s3         *s3.S3
	uploader   *s3manager.Uploader
	bucketName string
}

func (fs *S3FS) Open(path string) (io.ReadCloser, error) {
	params := &s3.GetObjectInput{
		Bucket: aws.String(fs.bucketName),
		Key:    aws.String(path),
	}

	resp, err := fs.s3.GetObject(params)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func (fs *S3FS) ReadAll(path string) ([]byte, error) {
	r, err := fs.Open(path)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(r)

	return buf.Bytes(), nil
}

func (fs *S3FS) Create(path string) (io.WriteCloser, error) {
	_, writer := io.Pipe()

	params := &s3manager.UploadInput{
		Bucket:       aws.String(fs.bucketName),
		Key:          aws.String(path),
		ACL:          aws.String(s3.ObjectCannedACLBucketOwnerRead),
		StorageClass: aws.String(s3.StorageClassStandardIa),
		Body:         strings.NewReader("Hello!"),
		// Body:         reader,
	}

	_, err := fs.uploader.Upload(params)
	if err != nil {
		return nil, err
	}

	return writer, nil
}

func (fs *S3FS) WriteAll(path string, contents []byte) error {
	w, err := fs.Create(path)
	if err != nil {
		return err
	}
	defer w.Close()

	_, err = io.Copy(w, bytes.NewReader(contents))

	return err
}

func (fs *S3FS) Remove(path string) error {
	params := &s3.DeleteObjectInput{
		Bucket: aws.String(fs.bucketName),
		Key:    aws.String(path),
	}

	_, err := fs.s3.DeleteObject(params)

	return err
}

// NewS3FS creates a new S3FS where all
// files are opened and created relative to the
func NewS3FS() *S3FS {

	bucketName := getBucketName()

	service, uploader := createS3Service()
	return &S3FS{
		s3:         service,
		uploader:   uploader,
		bucketName: bucketName,
	}
}

func getBucketName() string {
	bucketName := os.Getenv("AWS_BUCKET_NAME")
	if bucketName == "" {
		panic("Please enter your AWS_BUCKET_NAME environment variable!")
	}
	return bucketName
}

func createS3Service() (*s3.S3, *s3manager.Uploader) {
	sess := session.Must(session.NewSession())
	// TODO verify defaul region
	config := aws.NewConfig().WithRegion("us-west-2").WithLogLevel(aws.LogDebug)
	service := s3.New(sess, config)

	sess, err := session.NewSession(config)
	if err != nil {
		panic(err)
	}

	uploader := s3manager.NewUploader(sess)
	return service, uploader
}
