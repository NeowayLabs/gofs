package gofs

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3 struct {
	s3Client   *s3.S3
	uploader   *s3manager.Uploader
	bucketName string
}

func (fs *S3) listObjects(path string) (*s3.ListObjectsV2Output, error) {
	handledPath := path

	if strings.HasPrefix(handledPath, "/") {
		handledPath = strings.Replace(handledPath, "/", "", 1)
	}

	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(fs.bucketName),
		Prefix: aws.String(handledPath),
		// Prefix: aws.String("gofs"),
		// MaxKeys: aws.Int64(100),
	}

	return fs.s3Client.ListObjectsV2(input)
}

func (fs *S3) Open(path string) (io.ReadCloser, error) {
	fs.listObjects(path)

	input := &s3.GetObjectInput{
		Bucket: aws.String(fs.bucketName),
		Key:    aws.String(path),
	}

	resp, err := fs.s3Client.GetObject(input)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func (fs *S3) ReadAll(path string) ([]byte, error) {
	r, err := fs.Open(path)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(r)

	return buf.Bytes(), nil
}

func (fs *S3) Create(path string) (io.WriteCloser, error) {
	reader, writer := io.Pipe()

	input := &s3manager.UploadInput{
		Bucket:       aws.String(fs.bucketName),
		Key:          aws.String(path),
		ACL:          aws.String(s3.ObjectCannedACLBucketOwnerRead),
		StorageClass: aws.String(s3.StorageClassStandardIa),
		Body:         reader,
	}

	go func() {
		_, err := fs.uploader.Upload(input)
		if err != nil {
			log.Fatalf("Could not upload file: %s", err)
		}
		defer reader.Close()
	}()

	return writer, nil
}

func (fs *S3) WriteAll(path string, contents []byte) error {
	w, err := fs.Create(path)
	if err != nil {
		return err
	}
	defer w.Close()

	_, err = io.Copy(w, bytes.NewReader(contents))

	return err
}

func (fs *S3) Remove(path string) error {
	result, err := fs.listObjects(path)
	if err != nil {
		return err
	}

	count := int(*result.KeyCount)

	if count == 0 {
		return fmt.Errorf("file '%s' does not exists", path)
	}

	for count > 0 {
		// Create Delete object with slots for the objects to delete
		var items s3.Delete
		var objs = make([]*s3.ObjectIdentifier, int(*result.KeyCount))

		for i, o := range result.Contents {
			// Add objects from command line to array
			objs[i] = &s3.ObjectIdentifier{Key: aws.String(*o.Key)}
		}

		// Add list of objects to delete to Delete object
		items.SetObjects(objs)

		// Delete the items
		_, err = fs.s3Client.DeleteObjects(
			&s3.DeleteObjectsInput{Bucket: &fs.bucketName, Delete: &items},
		)
		if err != nil {
			return err
		}

		result, err = fs.listObjects(path)
		if err != nil {
			return err
		}
		count = int(*result.KeyCount)
	}

	return nil
}

// NewS3 creates a new S3 where all
// files are opened and created relative to the bucket name
func NewS3() *S3 {

	bucketName, err := getBucketName()
	if err != nil {
		log.Fatalf("Could not get bucket name: %s", err)
	}

	s3, err := createS3FS(bucketName)
	if err != nil {
		log.Fatalf("Could not create s3 file system: %s", err)
	}

	return s3
}

func getBucketName() (string, error) {
	bucketName := os.Getenv("AWS_BUCKET_NAME")
	if bucketName == "" {
		return "", fmt.Errorf("Please enter your AWS_BUCKET_NAME environment variable!")
	}
	return bucketName, nil
}

func createS3FS(bucketName string) (*S3, error) {
	sess := session.Must(session.NewSession())

	// TODO verify defaul region
	config := aws.NewConfig().WithRegion("us-west-2").WithLogLevel(aws.LogDebug)
	service := s3.New(sess, config)

	sess, err := session.NewSession(config)
	if err != nil {
		return nil, err
	}

	uploader := s3manager.NewUploader(sess)

	return &S3{
		s3Client:   service,
		uploader:   uploader,
		bucketName: bucketName,
	}, nil
}
