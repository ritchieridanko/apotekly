package cache

import (
	"context"
	"fmt"

	"github.com/ritchieridanko/apotekly/services/auth/internal/constants"
	cc "github.com/ritchieridanko/apotekly/services/shared/infra/cache"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
)

type AuthCache interface {
	UnreserveEmail(ctx context.Context, email string) (err *ce.Error)
	IsEmailReserved(ctx context.Context, email string) (reserved bool, err *ce.Error)
}

type authCache struct {
	cache *cc.Cache
}

func NewAuthCache(cc *cc.Cache) AuthCache {
	return &authCache{cache: cc}
}

func (c *authCache) UnreserveEmail(ctx context.Context, email string) *ce.Error {
	err := c.cache.Delete(
		ctx,
		constants.CachePrefixEmailReservation+":"+email,
	)
	if err != nil {
		return ce.NewError(
			ce.CodeCacheCommandExec,
			ce.MsgInternalServer,
			fmt.Errorf("failed to unreserve email: %w", err),
		)
	}
	return nil
}

func (c *authCache) IsEmailReserved(ctx context.Context, email string) (bool, *ce.Error) {
	exists, err := c.cache.Exists(
		ctx,
		constants.CachePrefixEmailReservation+":"+email,
	)
	if err != nil {
		return false, ce.NewError(
			ce.CodeCacheCommandExec,
			ce.MsgInternalServer,
			fmt.Errorf("failed to check if email is reserved: %w", err),
		)
	}
	return exists, nil
}
