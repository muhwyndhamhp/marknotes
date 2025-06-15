package cloudbucket

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	goconf "github.com/muhwyndhamhp/marknotes/config"
	"github.com/muhwyndhamhp/marknotes/utils/imageprocessing"
	"github.com/muhwyndhamhp/marknotes/utils/storage"
	"io"
	"log"
	"mime/multipart"
	"os"
	"strings"
)

type S3Client struct {
	client *s3.Client
	iproc  *imageprocessing.Client
}

const defaultBucketName = "mwyndham-dev"
const dbBucket = "db-bucket"

func NewS3Client(iproc *imageprocessing.Client) *S3Client {
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

	return &S3Client{client, iproc}
}
func (c *S3Client) UploadDB(ctx context.Context, dbName string) error {
	file, err := os.Open(dbName)
	if err != nil {
		return err
	}

	_, err = c.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(dbBucket),
		Key:    aws.String(dbName),
		Body:   file,
		Metadata: map[string]string{
			"Cache-Control": "no-store, no-cache, must-revalidate",
			"Expires":       "0",
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *S3Client) DownloadDB(ctx context.Context, dbName string) error {
	f, err := os.Create(dbName)
	if err != nil {
		return err
	}
	defer f.Close()

	out, err := c.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket:               aws.String(dbBucket),
		Key:                  aws.String(dbName),
		ResponseCacheControl: aws.String("no-store, no-cache, must-revalidate"),
	})
	if err != nil {
		return err
	}

	defer out.Body.Close()

	_, err = io.Copy(f, out.Body)
	if err != nil {
		return err
	}

	return nil
}

func ensureDir(path string) error {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		// Directory doesn't exist, create it
		return os.MkdirAll(path, 0755)
	} else if err != nil {
		// Some other error accessing the path
		return err
	}

	if !info.IsDir() {
		return fmt.Errorf("%s exists but is not a directory", path)
	}

	return nil
}

func (c *S3Client) UploadStatic(ctx context.Context, filename, exludePrefix, contentType string, bucketName string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}

	key := file.Name()
	if exludePrefix != "" {
		// remove the prefix from the key using match substring
		key = strings.SplitAfter(key, exludePrefix)[1]
	}

	if bucketName == "" {
		bucketName = defaultBucketName
	}
	_, err = c.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("https://resource.mwyndham.dev/%s", file.Name()), nil
}

func (c *S3Client) UploadMedia(ctx context.Context, f *multipart.FileHeader, prefix string, contentType string, size int) (string, error) {
	fname := strings.ReplaceAll(f.Filename, " ", "_")
	name := ""
	obj := &s3.PutObjectInput{
		Bucket:      aws.String(defaultBucketName),
		ContentType: aws.String(contentType),
	}

	if contentType != "image/gif" {
		r, size, err := c.iproc.ResizeImage(f, size)
		if err != nil {
			return "", err
		}
		intSize := int64(size)

		name = fmt.Sprintf("%s-%s", prefix, storage.AppendTimestamp(fname, ".webp"))
		obj.Key = aws.String(name)
		obj.Body = r
		obj.ContentLength = &intSize
		obj.ContentType = aws.String("image/webp")
	} else {
		name = fmt.Sprintf("%s-%s", prefix, storage.AppendTimestamp(fname, ""))
		obj.Key = aws.String(name)
		file, err := f.Open()
		if err != nil {
			return "", err
		}
		obj.Body = file
	}
	_, err := c.client.PutObject(ctx, obj)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("https://resource.mwyndham.dev/%s", name), nil
}
