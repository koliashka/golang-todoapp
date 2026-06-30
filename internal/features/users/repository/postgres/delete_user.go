package users_postgres_repository

import (
	"context"
	"fmt"
)

func (r *UsersRepository) DeleteUser(
	ctx context.Context,
	id int,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	DELETE FROM todoapp.users
	WHERE id=$1;
	`

	cmdTag, err := r.pool.Exec(ctx, query, id)

	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("user with id='%d': %w", err)
	}

	return nil
}
