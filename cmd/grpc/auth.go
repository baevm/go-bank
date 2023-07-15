package grpc

import (
	"context"
	"fmt"
	"go-bank/internal/token"
	"strings"

	"google.golang.org/grpc/metadata"
)

const (
	authorizationTokenKey = "authorization"
	authorizationTypeKey  = "bearer"
)

func (s *GrpcServer) authorize(ctx context.Context) (*token.Payload, error) {
	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}

	vals := md.Get(authorizationTokenKey)

	if len(vals) == 0 {
		return nil, fmt.Errorf("missing authorization header")
	}

	authHeader := vals[0]

	fields := strings.Fields(authHeader)

	if len(fields) < 2 {
		return nil, fmt.Errorf("authorization header incorrect format")
	}

	if strings.ToLower(fields[0]) != authorizationTypeKey {
		return nil, fmt.Errorf("authorization header incorrect type")
	}

	token := fields[1]
	payload, err := s.tokenMaker.Verify(token)

	if err != nil {
		return nil, fmt.Errorf("invalid authorization token")
	}

	return payload, nil
}
