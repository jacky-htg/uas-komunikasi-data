package storage

import (
	"context"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Connection() (*minio.Client, error) {
	return minio.New(os.Getenv("MINIO_HOST")+":"+os.Getenv("MINIO_PORT"), &minio.Options{
		Creds:  credentials.NewStaticV4(os.Getenv("MINIO_ACCESS_KEY"), os.Getenv("MINIO_SECRET_KEY"), ""),
		Secure: false, // Ganti menjadi true jika menggunakan HTTPS
	})
}

func SignUrl(ctx context.Context, minioClient *minio.Client, object string) (string, error) {
	expiryDuration, err := time.ParseDuration(os.Getenv("MINIO_EXPIRES"))
	if err != nil {
		return "", status.Errorf(codes.Internal, "Gagal parsing duration dari environment variable: %v", err)
	}

	presignedURL, err := minioClient.PresignedGetObject(ctx, os.Getenv("MINIO_BUCKET"), object, expiryDuration, nil)
	if err != nil {
		return "", status.Errorf(codes.Internal, "Gagal PresignedGetObject: %v", err)
	}

	return presignedURL.String(), nil
}
