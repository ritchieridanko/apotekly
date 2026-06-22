package cache

import (
	"context"
	"fmt"
	"strconv"

	"github.com/ritchieridanko/apotekly/services/auth/internal/constants"
	"github.com/ritchieridanko/apotekly/services/auth/internal/models"
	cc "github.com/ritchieridanko/apotekly/services/shared/infra/cache"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
)

type TokenCache interface {
	CreateVerification(ctx context.Context, data *models.CreateVerificationToken) (err *ce.Error)
}

type tokenCache struct {
	cache *cc.Cache
}

func NewTokenCache(cc *cc.Cache) TokenCache {
	return &tokenCache{cache: cc}
}

func (c *tokenCache) CreateVerification(ctx context.Context, data *models.CreateVerificationToken) *ce.Error {
	prefix := constants.CachePrefixEmailVerification
	script := `
		local token = redis.call("GET", KEYS[1])
		if token then
			redis.call("DEL", KEYS[1])
			redis.call("DEL", KEYS[3] .. ":" .. token)
		end

		redis.call("SET", KEYS[1], ARGV[1], "EX", ARGV[3])
		redis.call("SET", KEYS[2], ARGV[2], "EX", ARGV[3])
		return 1
	`

	_, err := c.cache.Evaluate(
		ctx, "s:crever", script,
		[]string{
			prefix + ":" + strconv.FormatUint(data.AuthID, 10),
			prefix + ":" + data.Token,
			prefix,
		},
		[]any{
			data.Token,
			data.AuthID,
			int(data.Duration.Seconds()),
		},
	)
	if err != nil {
		return ce.NewError(
			ce.CodeCacheScriptExec,
			ce.MsgInternalServer,
			fmt.Errorf("failed to create verification token: %w", err),
		)
	}

	return nil
}
