package cache

import (
	"context"
	"errors"
	"fmt"

	"github.com/ritchieridanko/apotekly/services/auth/internal/constants"
	"github.com/ritchieridanko/apotekly/services/auth/internal/models"
	cc "github.com/ritchieridanko/apotekly/services/shared/infra/cache"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
)

type OAuthCache interface {
	StoreState(ctx context.Context, data *models.StoreOAuthState) (err *ce.Error)
	UseState(ctx context.Context, state string) (provider string, err *ce.Error)
}

type oAuthCache struct {
	cache *cc.Cache
}

func NewOAuthCache(cc *cc.Cache) OAuthCache {
	return &oAuthCache{cache: cc}
}

func (c *oAuthCache) StoreState(ctx context.Context, data *models.StoreOAuthState) *ce.Error {
	err := c.cache.Set(
		ctx,
		constants.CachePrefixOAuthState+":"+data.State,
		data.Provider,
		data.Duration,
	)
	if err != nil {
		return ce.NewError(
			ce.CodeCacheCommandExec,
			ce.MsgInternalServer,
			fmt.Errorf("failed to store oauth state: %w", err),
		)
	}
	return nil
}

func (c *oAuthCache) UseState(ctx context.Context, state string) (string, *ce.Error) {
	script := `
		local provider = redis.call("GET", KEYS[1])
		if provider then
			redis.call("DEL", KEYS[1])
			return provider
		end
		return nil
	`

	res, err := c.cache.Evaluate(
		ctx, "s:usoast", script,
		[]string{
			constants.CachePrefixOAuthState + ":" + state, // KEYS[1]
		},
		nil,
	)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to use oauth state: %w", err)
		if errors.Is(err, ce.ErrCacheNoResult) {
			return "", ce.NewError(
				ce.CodeInvalidOAuthState,
				ce.MsgInvalidOAuthState,
				wrappedErr,
			)
		}
		return "", ce.NewError(
			ce.CodeCacheScriptExec,
			ce.MsgInternalServer,
			wrappedErr,
		)
	}

	value, ok := res.(string)
	if !ok {
		return "", ce.NewError(ce.CodeTypeAssertionFailed, ce.MsgInternalServer, nil)
	}

	return value, nil
}
