package grpc

import (
	"context"
	"database/sql"
	"errors"
	db "go-bank/db/sqlc"
	"go-bank/internal/password"
	"go-bank/pb"
	"time"

	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *GrpcServer) CreateUser(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	hashedPass, err := password.Hash(req.Password)

	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	args := db.CreateUserParams{
		Email:      req.Email,
		HashedPass: hashedPass,
		Username:   req.Username,
		FullName:   req.FullName,
	}

	_, err = s.db.CreateUser(ctx, args)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "username already exists")
			}

			return nil, status.Errorf(codes.Internal, "failed to create user")
		}
	}

	return &pb.CreateResponse{
		Message: "ok",
	}, nil
}

func (s *GrpcServer) LoginUser(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := s.db.GetUser(ctx, req.Username)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		} else {
			return nil, status.Errorf(codes.Internal, "internal server error")
		}
	}

	err = password.Check(user.HashedPass, req.Password)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "incorrect password")
	}

	accessToken, accessPayload, err := s.tokenMaker.Create(req.Username, s.cfg.ACCESS_TOKEN_DURATION)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session")
	}

	refreshToken, refreshPayload, err := s.tokenMaker.Create(req.Username, s.cfg.REFRESH_TOKEN_DURATION)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session")
	}

	metadata := s.extractMetadata(ctx)

	session, err := s.db.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		IsBlocked:    false,
		UserAgent:    metadata.UserAgent,
		ClientIp:     metadata.ClientIp,
		ExpiresAt:    refreshPayload.ExpiresAt.Time,
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session")
	}

	return &pb.LoginResponse{
		SessionId:          session.ID.String(),
		Username:           user.Username,
		AccessToken:        accessToken,
		AccessTokenExpire:  timestamppb.New(accessPayload.ExpiresAt.Time),
		RefreshToken:       refreshToken,
		RefreshTokenExpire: timestamppb.New(refreshPayload.ExpiresAt.Time),
	}, nil
}

func (s *GrpcServer) UpdateUser(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {

	args := db.UpdateUserParams{
		Username: req.Username,
		FullName: sql.NullString{
			String: req.GetFullName(),
			Valid:  req.FullName != nil,
		},
		Email: sql.NullString{
			String: req.GetEmail(),
			Valid:  req.Email != nil,
		},
	}

	if req.Password != nil {
		hashedPass, err := password.Hash(req.GetPassword())

		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}

		args.HashedPass = sql.NullString{
			String: hashedPass,
			Valid:  true,
		}

		args.PasswordChangedAt = sql.NullTime{
			Time: time.Now(),
			Valid: true,
		}
	}

	_, err := s.db.UpdateUser(ctx, args)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows){
			return nil, status.Errorf(codes.NotFound, "user not found")
		} else {
			return nil, status.Errorf(codes.Internal, "failed to update user")
		}
	}

	return &pb.UpdateResponse{
		Message: "ok",
	}, nil
}
