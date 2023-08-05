package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	db "go-bank/db/sqlc"
	"go-bank/internal/testutil"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

type VerifyEmailPayload struct {
	Username string
}

const (
	TypeEmailVerify = "email:send_verify"
)

func (d *RedisTaskDistributor) SendVerifyEmailTask(ctx context.Context, payload *VerifyEmailPayload, opts ...asynq.Option) error {
	jsonPayload, err := json.Marshal(payload)

	if err != nil {
		return err
	}

	task := asynq.NewTask(TypeEmailVerify, jsonPayload, opts...)

	info, err := d.client.EnqueueContext(ctx, task)

	if err != nil {
		return err
	}

	log.Info().
		Str("type", info.Type).
		Str("queue", info.Queue).
		Str("task id", info.ID).
		Msg("enqueued task")

	return err
}

func (p *RedisTaskProcessor) ProcessSendVerifyEmailTask(ctx context.Context, t *asynq.Task) error {
	var payload VerifyEmailPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	user, err := p.store.GetUser(ctx, payload.Username)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("user doesnt exist: %w", asynq.SkipRetry)
		} else {
			return fmt.Errorf("failed to get user: %w", err)
		}
	}

	verifyEmail, err := p.store.CreateVerifyEmail(ctx, db.CreateVerifyEmailParams{
		Username:   user.Username,
		Email:      user.Email,
		SecretCode: testutil.RandomString(32),
	})

	if err != nil {
		return err
	}

	verifyLink := fmt.Sprintf("http://localhost:5000/v1/verify_email?email_id=%d&secret_code=%s", verifyEmail.ID, verifyEmail.SecretCode)

	subj := "Confirm your email"
	content := fmt.Sprintf(`
	<h1>Confirm your email by clicking this link</h1>
	<div>
	<a href="%s">Click</a>
	</div>
	`, verifyLink)

	err = p.mailer.SendEmail(subj, content, []string{user.Email}, nil, nil, nil)

	if err != nil {
		return err
	}

	log.Info().
		Str("type", t.Type()).
		Str("user_email", user.Email).
		Msg("processed task")

	return nil
}
