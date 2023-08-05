package worker

import (
	"context"

	"github.com/hibiken/asynq"
)

type TaskDistributor interface {
	SendVerifyEmailTask(ctx context.Context, payload *VerifyEmailPayload, opts ...asynq.Option) error
}

type RedisTaskDistributor struct {
	client *asynq.Client
}

func NewRedisTaskDistributor(redis asynq.RedisClientOpt) TaskDistributor {
	client := asynq.NewClient(redis)

	return &RedisTaskDistributor{
		client: client,
	}
}
