package db

import (
	"context"
	"database/sql"
	"fmt"
)

type VerifyEmailTxParams struct {
	EmailId    int64
	SecretCode string
}

type VerifyEmailTxResult struct {
	VerifyEmail VerifyEmails
	User        Users
}

func (s *SQLStore) VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error) {
	var res VerifyEmailTxResult

	err := s.execTx(ctx, func(q *Queries) error {
		var err error

		res.VerifyEmail, err = q.UpdateVerifyEmail(ctx, UpdateVerifyEmailParams{
			ID:         arg.EmailId,
			SecretCode: arg.SecretCode,
		})

		if err != nil {
			return fmt.Errorf("verify email not found")
		}

		res.User, err = q.UpdateUser(ctx, UpdateUserParams{
			Username: res.VerifyEmail.Username,
			IsEmailActivated: sql.NullBool{
				Bool:  true,
				Valid: true,
			},
		})

		return err
	})

	return res, err
}
