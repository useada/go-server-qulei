package main

import (
	"bytes"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"a.com/go-server/proto/pb"
)

type S3Client struct {
	Client *s3.S3
	ACL    string // s3 acl for uploaded files - for our use either "public" or "private"
	Bucket string // s3 bucket to upload to
}

func (s *S3Client) Upload(id, ex string, typef pb.FileType, data []byte) error {
	_, err := s.Client.PutObject(&s3.PutObjectInput{
		Key:         aws.String(s.RemoteDir(id, ex, typef)),
		Body:        bytes.NewReader(data),
		Bucket:      aws.String(s.Bucket),
		ACL:         aws.String(s.ACL),
		ContentType: aws.String(ContentTypeMap[ex]),
	})
	return err
}

func (s *S3Client) Delete(remoteDir string) error {
	_, err := s.Client.DeleteObject(&s3.DeleteObjectInput{
		Key:    aws.String(remoteDir),
		Bucket: aws.String(s.Bucket),
	})
	return err
}

func (s *S3Client) GetBucket() string {
	return s.Bucket
}

var RemoteDirMap = map[pb.FileType]string{
	pb.FileType_IMAGE:  "/image/",
	pb.FileType_AVATAR: "/avatar/",
	pb.FileType_FILE:   "/file/",
}

func (s *S3Client) RemoteDir(id, ex string, typef pb.FileType) string {
	return RemoteDirMap[typef] + "original/" + id + "." + ex
}

var ContentTypeMap = map[string]string{
	"jpg":  "image/jpeg",
	"jpeg": "image/jpeg",
	"png":  "image/png",
	"gif":  "image/gif",
}

var S3 *S3Client

func InitS3Client(conf S3Configor) {
	S3 = &S3Client{
		Bucket: conf.Bucket,
		ACL:    conf.ACL,
		Client: s3.New(session.Must(session.NewSession()), &aws.Config{Region: aws.String(conf.Region)}),
	}
}