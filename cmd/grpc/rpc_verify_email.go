package grpc

import (
	"context"
	db "go-bank/db/sqlc"
	"go-bank/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *GrpcServer) VerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {
	_, err := s.db.VerifyEmailTx(ctx, db.VerifyEmailTxParams{
		EmailId:    req.GetEmailId(),
		SecretCode: req.GetSecretCode(),
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to verify email")
	}

	return &pb.VerifyEmailResponse{
		IsVerified: true,
	}, nil
}
