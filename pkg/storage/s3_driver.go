package storage

import (
	"context"
	"fmt"
	"io"
	"mime"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type s3Storage struct {
	client     *s3.Client
	uploader   *manager.Uploader
	downloader *manager.Downloader
	bucket     string
	region     string
	endpoint   string
}

// S3Options represents S3 storage options
type S3Options struct {
	Bucket    string `mapstructure:"bucket"`
	Region    string `mapstructure:"region"`
	Endpoint  string `mapstructure:"endpoint"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
}

func init() {
	Register("s3", NewS3Storage)
}

// NewS3Storage creates a new S3 storage instance
func NewS3Storage(config *Config) (Storage, error) {
	var opts S3Options
	if err := mapstructureDecodeConfig(config.Options, &opts); err != nil {
		return nil, err
	}

	// Create AWS config
	cfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion(opts.Region),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			opts.AccessKey,
			opts.SecretKey,
			"",
		)),
	)
	if err != nil {
		return nil, err
	}

	// Create S3 client
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		if opts.Endpoint != "" {
			o.BaseEndpoint = aws.String(opts.Endpoint)
		}
	})

	return &s3Storage{
		client:     client,
		uploader:   manager.NewUploader(client),
		downloader: manager.NewDownloader(client),
		bucket:     opts.Bucket,
		region:     opts.Region,
		endpoint:   opts.Endpoint,
	}, nil
}

func (s *s3Storage) Upload(ctx context.Context, path string, reader io.Reader, options ...UploadOption) error {
	if err := validatePath(path); err != nil {
		return err
	}

	// Apply options
	opts := &uploadOptions{}
	for _, opt := range options {
		opt(opts)
	}

	// Prepare upload input
	input := &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
		Body:   reader,
	}

	// Set content type if provided
	if opts.ContentType != "" {
		input.ContentType = aws.String(opts.ContentType)
	} else {
		contentType := mime.TypeByExtension(filepath.Ext(path))
		if contentType != "" {
			input.ContentType = aws.String(contentType)
		}
	}

	// Set ACL if public
	if opts.Public {
		input.ACL = types.ObjectCannedACLPublicRead
	}

	// Set cache control if max age is provided
	if opts.MaxAge > 0 {
		input.CacheControl = aws.String(fmt.Sprintf("max-age=%d", int(opts.MaxAge.Seconds())))
	}

	_, err := s.uploader.Upload(ctx, input)
	return err
}

func (s *s3Storage) Download(ctx context.Context, path string) (io.ReadCloser, error) {
	if err := validatePath(path); err != nil {
		return nil, err
	}

	// Create a buffer to write the file to
	buf := manager.NewWriteAtBuffer([]byte{})

	// Download the file
	_, err := s.downloader.Download(ctx, buf, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return nil, err
	}

	// Create a reader from the buffer
	return io.NopCloser(strings.NewReader(string(buf.Bytes()))), nil
}

func (s *s3Storage) Delete(ctx context.Context, path string) error {
	if err := validatePath(path); err != nil {
		return err
	}

	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	})
	return err
}

func (s *s3Storage) Exists(ctx context.Context, path string) (bool, error) {
	if err := validatePath(path); err != nil {
		return false, err
	}

	_, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return false, nil
	}
	return true, nil
}

func (s *s3Storage) List(ctx context.Context, prefix string) ([]FileInfo, error) {
	if err := validatePath(prefix); err != nil {
		return nil, err
	}

	var files []FileInfo
	paginator := s3.NewListObjectsV2Paginator(s.client, &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucket),
		Prefix: aws.String(prefix),
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, obj := range page.Contents {
			files = append(files, FileInfo{
				Name:         filepath.Base(aws.ToString(obj.Key)),
				Size:         aws.ToInt64(obj.Size),
				LastModified: aws.ToTime(obj.LastModified),
				IsDir:        strings.HasSuffix(aws.ToString(obj.Key), "/"),
				ContentType:  mime.TypeByExtension(filepath.Ext(aws.ToString(obj.Key))),
				Path:         aws.ToString(obj.Key),
			})
		}
	}

	return files, nil
}

func (s *s3Storage) GetURL(ctx context.Context, path string) (string, error) {
	if err := validatePath(path); err != nil {
		return "", err
	}

	if s.endpoint != "" {
		return fmt.Sprintf("%s/%s/%s", s.endpoint, s.bucket, path), nil
	}

	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucket, s.region, path), nil
}

func (s *s3Storage) GetInfo(ctx context.Context, path string) (*FileInfo, error) {
	if err := validatePath(path); err != nil {
		return nil, err
	}

	output, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return nil, err
	}

	return &FileInfo{
		Name:         filepath.Base(path),
		Size:         aws.ToInt64(output.ContentLength),
		LastModified: aws.ToTime(output.LastModified),
		IsDir:        strings.HasSuffix(path, "/"),
		ContentType:  aws.ToString(output.ContentType),
		Path:         path,
	}, nil
}

func (s *s3Storage) Close() error {
	return nil
}
