package s3

import (
	"context"
	"fmt"
	"io"
	"joerx/minecraft-cli/internal/api/backup"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// S3Store implements the handler.ObjectStore interface and stores objects in AWS S3 or compatible
// storage services
type S3Store struct {
	sess   *session.Session
	bucket string
}

func NewStore(sess *session.Session, bucket string) *S3Store {
	return &S3Store{sess, bucket}
}

func (st *S3Store) Put(ctx context.Context, key string, r io.Reader) (backup.ObjectInfo, error) {
	uploader := s3manager.NewUploader(st.sess)
	oi := backup.ObjectInfo{}

	_, err := uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket: aws.String(st.bucket),
		Key:    aws.String(key),
		Body:   r,
	})

	if err != nil {
		return oi, err
	}

	oi.Location = fmt.Sprintf("s3://%s/%s", st.bucket, key)
	return oi, nil
}

func (st *S3Store) List(ctx context.Context) ([]backup.ObjectInfo, error) {
	svc := s3.New(st.sess)
	result, err := svc.ListObjects(&s3.ListObjectsInput{
		Bucket:  aws.String(st.bucket),
		MaxKeys: aws.Int64(10),
	})

	if err != nil {
		log.Println(err)
		return []backup.ObjectInfo{}, s3Error(err)
	}

	log.Printf("Found %d objects in s3://%s", len(result.Contents), st.bucket)

	objects := make([]backup.ObjectInfo, len(result.Contents))
	for i, r := range result.Contents {
		objects[i] = backup.ObjectInfo{
			Key:      *r.Key,
			Location: fmt.Sprintf("s3://%s/%s", st.bucket, *r.Key),
		}
	}

	return objects, nil
}

func (st *S3Store) Get(ctx context.Context, key string, w io.WriterAt) error {
	dl := s3manager.NewDownloader(st.sess)
	in := &s3.GetObjectInput{
		Bucket: aws.String(st.bucket),
		Key:    aws.String(key),
	}

	if _, err := dl.Download(w, in); err != nil {
		return s3Error(err)
	}

	return nil
}

func s3Error(err error) error {
	var msg string
	splits := strings.Split(strings.ReplaceAll(err.Error(), "\r\n", "\n"), "\n")

	if len(splits) >= 1 {
		msg = splits[0]
	} else {
		msg = "unknown error"
	}

	return fmt.Errorf("s3 error - %s", msg)
}
