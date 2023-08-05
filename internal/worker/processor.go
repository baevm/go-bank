package worker

import (
	"context"
	sqlc "go-bank/db/sqlc"
	"go-bank/internal/mail"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

type TaskProcessor interface {
	Start() error
	ProcessSendVerifyEmailTask(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	store  sqlc.Store
	mailer mail.EmailSender
}

func NewRedisTaskProcessor(r asynq.RedisClientOpt, store sqlc.Store, mailer mail.EmailSender) TaskProcessor {
	server := asynq.NewServer(r, asynq.Config{
		ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
			log.Error().
				Err(err).
				Str("type", task.Type()).
				Msg("task process failed")
		}),
	})

	return &RedisTaskProcessor{
		server: server,
		store:  store,
		mailer: mailer,
	}
}

func (p *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(TypeEmailVerify, p.ProcessSendVerifyEmailTask)

	return p.server.Run(mux)
}
