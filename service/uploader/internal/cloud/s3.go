package cloud

import (
	"bytes"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"a.com/go-server/proto/pb"
)

type s3cloud struct {
	*S3Client
}

func NewS3Repo(client *S3Client) Cloud {
	return &s3cloud{client}
}

func (s *s3cloud) Upload(id, ex string, xtype pb.TYPE, data []byte) error {
	_, err := s.Client.PutObject(&s3.PutObjectInput{
		Key:         aws.String(s.RemoteDir(id, ex, xtype)),
		Body:        bytes.NewReader(data),
		Bucket:      aws.String(s.Bucket),
		ACL:         aws.String(s.ACL),
		ContentType: aws.String(ContentTypeMap[ex]),
	})
	return err
}

func (s *s3cloud) Delete(remoteDir string) error {
	_, err := s.Client.DeleteObject(&s3.DeleteObjectInput{
		Key:    aws.String(remoteDir),
		Bucket: aws.String(s.Bucket),
	})
	return err
}

func (s *s3cloud) GetBucket() string {
	return s.Bucket
}

var RemoteDirMap = map[pb.TYPE]string{
	pb.TYPE_IMAGE:  "/image/",
	pb.TYPE_AVATAR: "/avatar/",
	pb.TYPE_FILE:   "/file/",
}

func (s *s3cloud) RemoteDir(id, ex string, xtype pb.TYPE) string {
	return RemoteDirMap[xtype] + "original/" + id + "." + ex
}

var ContentTypeMap = map[string]string{
	"jpg":  "image/jpeg",
	"jpeg": "image/jpeg",
	"png":  "image/png",
	"gif":  "image/gif",
}

type S3Client struct {
	Client *s3.S3
	ACL    string // s3 acl for uploaded files - for our use either "public" or "private"
	Bucket string // s3 bucket to upload to
}

type Config struct {
	Region string
	ACL    string
	Bucket string
}

func NewS3Client(conf Config) *S3Client {
	client := &S3Client{
		Bucket: conf.Bucket,
		ACL:    conf.ACL,
		Client: s3.New(session.Must(session.NewSession()), &aws.Config{Region: aws.String(conf.Region)}),
	}

	fmt.Println("")
	return client
}
