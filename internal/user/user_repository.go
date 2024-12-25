package user

import "context"

type UserRepository interface {
	GetCurrentUserId(ctx context.Context) (string, error)
}
