package grpc

import (
	"context"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GrpcLogger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (res interface{}, err error) {
	startTime := time.Now()

	res, err = handler(ctx, req)

	statusCode := codes.Unknown

	if st, ok := status.FromError(err); ok {
		statusCode = st.Code()
	}

	duration := time.Since(startTime)

	logger := log.Info()

	if err != nil {
		logger = log.Error().Err(err)
	}

	logger.
		Str("protocol", "grpc").
		Str("method", info.FullMethod).
		Dur("duration", duration).
		Str("status_text", statusCode.String()).
		Int("status_code", int(statusCode)).
		Msg("grpc request")

	return
}

type HttpRecorder struct {
	http.ResponseWriter
	StatusCode int
	Body       []byte
}

func (r *HttpRecorder) WriteHeader(status int) {
	r.StatusCode = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *HttpRecorder) Write(body []byte) (int, error) {
	r.Body = body
	return r.ResponseWriter.Write(body)
}

func GatewayLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		logger := log.Info()

		recorder := &HttpRecorder{
			ResponseWriter: w,
			StatusCode:     http.StatusOK,
		}

		handler.ServeHTTP(recorder, r)

		duration := time.Since(startTime)

		if recorder.StatusCode > 399 {
			logger = log.Error().Bytes("error", recorder.Body)
		}

		logger.
			Str("protocol", "http").
			Str("method", r.Method).
			Str("path", r.RequestURI).
			Dur("duration", duration).
			Str("status_text", http.StatusText(recorder.StatusCode)).
			Int("status_code", recorder.StatusCode).
			Msg("grpc gateway request")

	})
}
