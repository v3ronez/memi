package utils

import (
	"context"

	"github.com/v3ronez/memi/types"
)

func GetUserFromCtx(ctx context.Context) *types.User {
	user, ok := ctx.Value(types.UserCtxKey).(*types.User)
	if !ok {
		return nil
	}

	return user
}
