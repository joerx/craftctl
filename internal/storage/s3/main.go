package s3

import (
	"context"
	"fmt"
	"io"
	"joerx/minecraft-cli/internal/api/backup"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
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
