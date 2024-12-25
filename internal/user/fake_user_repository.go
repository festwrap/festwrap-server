package user

import "context"

type GetCurrentIdArgs struct {
	Context context.Context
}

type GetCurrentIdValue struct {
	UserId string
	Err    error
}

type FakeUserRepository struct {
	currentIdArgs  GetCurrentIdArgs
	currentIdValue GetCurrentIdValue
}

func (r *FakeUserRepository) GetCurrentUserId(ctx context.Context) (string, error) {
	r.currentIdArgs = GetCurrentIdArgs{Context: ctx}
	return r.currentIdValue.UserId, r.currentIdValue.Err

}

func (r FakeUserRepository) GetGetCurrentIdArgs() GetCurrentIdArgs {
	return r.currentIdArgs
}

func (r *FakeUserRepository) SetGetCurrentIdValue(value GetCurrentIdValue) {
	r.currentIdValue = value
}
