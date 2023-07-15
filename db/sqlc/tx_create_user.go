package db

import "context"

type CreateUserTxParams struct {
	CreateUserParams
	AfterCreate func(user Users) error
}

type CreateUserTxResult struct {
	User Users
}

func (s *SQLStore) CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error) {
	var res CreateUserTxResult

	err := s.execTx(ctx, func(q *Queries) error {
		var err error

		res.User, err = q.CreateUser(ctx, arg.CreateUserParams)

		if err != nil {
			return err
		}

		err = arg.AfterCreate(res.User)

		return err
	})

	return res, err
}
