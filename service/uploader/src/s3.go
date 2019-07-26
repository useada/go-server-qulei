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

func (s *S3Client) Upload(id, ex string, xtype pb.TYPE, data []byte) error {
	_, err := s.Client.PutObject(&s3.PutObjectInput{
		Key:         aws.String(s.RemoteDir(id, ex, xtype)),
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

var RemoteDirMap = map[pb.TYPE]string{
	pb.TYPE_IMAGE:  "/image/",
	pb.TYPE_AVATAR: "/avatar/",
	pb.TYPE_FILE:   "/file/",
}

func (s *S3Client) RemoteDir(id, ex string, xtype pb.TYPE) string {
	return RemoteDirMap[xtype] + "original/" + id + "." + ex
}

var ContentTypeMap = map[string]string{
	"jpg":  "image/jpeg",
	"jpeg": "image/jpeg",
	"png":  "image/png",
	"gif":  "image/gif",
}

var S3 *S3Client

func InitS3Client(conf S3Config) {
	S3 = &S3Client{
		Bucket: conf.Bucket,
		ACL:    conf.ACL,
		Client: s3.New(session.Must(session.NewSession()), &aws.Config{Region: aws.String(conf.Region)}),
	}
}