package grpc

import (
	"context"
	"fmt"
	"go-bank/config"
	db "go-bank/db/sqlc"
	"go-bank/doc"
	"go-bank/internal/token"
	"go-bank/pb"
	"go-bank/worker"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

type GrpcServer struct {
	db          db.Store
	tokenMaker  token.Maker
	cfg         config.Config
	distributor worker.TaskDistributor

	pb.UnimplementedUserServiceServer
}

func NewGrpcServer(config config.Config, db db.Store, distributor worker.TaskDistributor) (*GrpcServer, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TOKEN_SYMMETRIC_KEY)

	if err != nil {
		return nil, fmt.Errorf("cant create server: %s", err)
	}

	server := &GrpcServer{
		db:          db,
		tokenMaker:  tokenMaker,
		cfg:         config,
		distributor: distributor,
	}

	return server, nil
}

func (s *GrpcServer) Start() error {
	lis, err := net.Listen("tcp", s.cfg.GRPC_ADDR)

	if err != nil {
		return err
	}

	// Create grpc server with logger
	opts := grpc.UnaryInterceptor(GrpcLogger)
	grpcServer := grpc.NewServer(opts)
	pb.RegisterUserServiceServer(grpcServer, s)
	reflection.Register(grpcServer)

	// Start server
	log.Printf("Listening and serving GRPC on: %s \n", lis.Addr())
	err = grpcServer.Serve(lis)

	return err
}

func (s *GrpcServer) StartGateway() error {
	lis, err := net.Listen("tcp", s.cfg.SRV_ADDR)

	if err != nil {
		return err
	}

	// use snake_case as defined in proto file instead of camelCase
	protoNames := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	// Register grpc gateway
	grpcMux := runtime.NewServeMux(protoNames)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = pb.RegisterUserServiceHandlerServer(ctx, grpcMux, s)

	if err != nil {
		return err
	}

	// Create http mux
	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	// Initialize swagger
	swaggerFolder, err := doc.GetSwaggerFolder()

	if err != nil {
		return err
	}

	static := http.FileServer(http.FS(swaggerFolder))
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", static))

	// Handler with structured logger
	handler := GatewayLogger(mux)

	// Start server
	log.Printf("Listening and serving GRPC Gateway on: %s \n", s.cfg.SRV_ADDR)
	err = http.Serve(lis, handler)

	return err
}
