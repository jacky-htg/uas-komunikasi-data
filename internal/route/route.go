package route

import (
	"database/sql"
	"log"
	"uas-komdat/internal/domain"
	userPb "uas-komdat/pb/users"

	"github.com/minio/minio-go/v7"
	"google.golang.org/grpc"
)

// GrpcRoute func
func GrpcRoute(grpcServer *grpc.Server, db *sql.DB, log *log.Logger, minioCLient *minio.Client) {
	userServer := domain.UserService{Db: db, MinioClient: minioCLient}
	userPb.RegisterUserServiceServer(grpcServer, &userServer)
}
