package grpc

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

type Metadata struct {
	ClientIp  string
	UserAgent string
}

const (
	userAgentKey        = "user-agent"
	userAgentGatewayKey = "grpcgateway-user-agent"
	clientIpKey         = "x-forwarded-for"
)

func (s *GrpcServer) extractMetadata(ctx context.Context) *Metadata {
	m := &Metadata{}

	// get user agent from metadata
	md, ok := metadata.FromIncomingContext(ctx)

	if ok {
		// get user agent from grpc
		if userAgent := md.Get(userAgentKey); len(userAgent) > 0 {
			m.UserAgent = userAgent[0]
		}

		// get user agent from grpc gateway
		if userAgent := md.Get(userAgentGatewayKey); len(userAgent) > 0 {
			m.UserAgent = userAgent[0]
		}

		// get user ip from grpc gateway
		clientIp := md.Get(clientIpKey)

		if len(clientIp) > 0 {
			m.ClientIp = clientIp[0]
		}
	}

	// get user ip from peer from grpc
	if p, ok := peer.FromContext(ctx); ok {
		m.ClientIp = p.Addr.String()
	}

	return m
}
