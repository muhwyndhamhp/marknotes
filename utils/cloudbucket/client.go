package cloudbucket

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	goconf "github.com/muhwyndhamhp/marknotes/config"
	"github.com/muhwyndhamhp/marknotes/utils/storage"
)

type S3Client struct {
	client *s3.Client
}

const defaultBucketName = "mwyndham-dev"

func NewS3Client() *S3Client {
	accountId := goconf.Get(goconf.CF_ACCOUNT_ID)
	accessKeyId := goconf.Get(goconf.CF_R2_ACCESS_KEY_ID)
	accessKeySecret := goconf.Get(goconf.CF_R2_SECRET_ACCESS_KEY)

	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountId),
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyId, accessKeySecret, "")),
		config.WithRegion("auto"),
	)
	if err != nil {
		log.Fatal(err)
	}

	client := s3.NewFromConfig(cfg)

	return &S3Client{client}
}

func (c *S3Client) UploadStatic(ctx context.Context, filename string, contentType string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}

	_, err = c.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(defaultBucketName),
		Key:         aws.String(file.Name()),
		Body:        file,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("https://resource.mwyndham.dev/%s", file.Name()), nil
}

func (c *S3Client) UploadMedia(ctx context.Context, f *multipart.FileHeader, prefix string, contentType string) (string, error) {
	file, err := f.Open()
	if err != nil {
		return "", err
	}

	defer file.Close()

	fname := strings.ReplaceAll(f.Filename, " ", "_")
	name := fmt.Sprintf("%s-%s", prefix, storage.AppendTimestamp(fname))

	_, err = c.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(defaultBucketName),
		Key:         aws.String(name),
		Body:        file,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("https://resource.mwyndham.dev/%s", name), nil
}
