package cache

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/ritchieridanko/apotekly/services/auth/internal/constants"
	"github.com/ritchieridanko/apotekly/services/auth/internal/models"
	cc "github.com/ritchieridanko/apotekly/services/shared/infra/cache"
	"github.com/ritchieridanko/apotekly/services/shared/utils"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
)

type TokenCache interface {
	CreateEmailChange(ctx context.Context, data *models.CreateEmailChangeToken) (err *ce.Error)
	UseEmailChange(ctx context.Context, token string) (authID uint64, newEmail string, err *ce.Error)
	CreatePasswordReset(ctx context.Context, data *models.CreatePasswordResetToken) (err *ce.Error)
	CreateVerification(ctx context.Context, data *models.CreateVerificationToken) (err *ce.Error)
	UseVerification(ctx context.Context, token string) (authID uint64, err *ce.Error)
}

type tokenCache struct {
	cache *cc.Cache
}

func NewTokenCache(cc *cc.Cache) TokenCache {
	return &tokenCache{cache: cc}
}

func (c *tokenCache) CreateEmailChange(ctx context.Context, data *models.CreateEmailChangeToken) *ce.Error {
	emch := constants.CachePrefixEmailChange
	emres := constants.CachePrefixEmailReservation
	script := `
		local token = redis.call("GET", KEYS[1])
		if token then
			local email = redis.call("HGET", KEYS[3] .. ":" .. token, "ne")
			if email then
				redis.call("DEL", KEYS[5] .. ":" .. email)
			end

			redis.call("DEL", KEYS[1])
			redis.call("DEL", KEYS[3] .. ":" .. token)
		end

		local reserved = redis.call("SET", KEYS[4], ARGV[2], "NX", "EX", ARGV[4])
		if not reserved then
			return 0
		end

		redis.call("SET", KEYS[1], ARGV[1], "EX", ARGV[4])
		redis.call("HSET", KEYS[2], "id", ARGV[2], "ne", ARGV[3])
		redis.call("EXPIRE", KEYS[2], ARGV[4])
		return 1
	`

	res, err := c.cache.Evaluate(
		ctx, "s:cremch", script,
		[]string{
			emch + ":" + strconv.FormatUint(data.AuthID, 10), // KEYS[1]
			emch + ":" + data.Token,                          // KEYS[2]
			emch,                                             // KEYS[3]
			emres + ":" + data.NewEmail,                      // KEYS[4]
			emres,                                            // KEYS[5]
		},
		[]any{
			data.Token,                   // ARGV[1]
			data.AuthID,                  // ARGV[2]
			data.NewEmail,                // ARGV[3]
			int(data.Duration.Seconds()), // ARGV[4]
		},
	)
	if err != nil {
		return ce.NewError(
			ce.CodeCacheScriptExec,
			ce.MsgInternalServer,
			fmt.Errorf("failed to create email change token: %w", err),
		)
	}
	if res.(int64) == 0 {
		return ce.NewError(ce.CodeEmailNotAvailable, ce.MsgEmailAlreadyRegistered, nil)
	}

	return nil
}

func (c *tokenCache) UseEmailChange(ctx context.Context, token string) (uint64, string, *ce.Error) {
	emch := constants.CachePrefixEmailChange
	script := `
		if redis.call("EXISTS", KEYS[1]) == 0 then
			return nil
		end

		local data = redis.call("HMGET", KEYS[1], "id", "ne")
		local authID = data[1]
		local email = data[2]
		
		redis.call("DEL", KEYS[1])
		redis.call("DEL", KEYS[2] .. ":" .. authID)
		
		return {authID, email}
	`

	res, err := c.cache.Evaluate(
		ctx, "s:usemch", script,
		[]string{
			emch + ":" + token, // KEYS[1]
			emch,               // KEYS[2]
		},
		nil,
	)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to use email change token: %w", err)
		if errors.Is(err, ce.ErrCacheNoResult) {
			return 0, "", ce.NewError(
				ce.CodeInvalidToken,
				ce.MsgInvalidToken,
				wrappedErr,
			)
		}
		return 0, "", ce.NewError(
			ce.CodeCacheScriptExec,
			ce.MsgInternalServer,
			wrappedErr,
		)
	}

	values, ok := res.([]any)
	if !ok || len(values) != 2 {
		return 0, "", ce.NewError(ce.CodeTypeAssertionFailed, ce.MsgInternalServer, nil)
	}

	id, ok1 := values[0].(string)
	newEmail, ok2 := values[1].(string)
	if !ok1 || !ok2 {
		return 0, "", ce.NewError(ce.CodeTypeAssertionFailed, ce.MsgInternalServer, nil)
	}

	authID, err := utils.ToUint64(id)
	if err != nil {
		return 0, "", ce.NewError(
			ce.CodeTypeConversionFailed,
			ce.MsgInternalServer,
			fmt.Errorf("failed to use email change token: %w", err),
		)
	}

	return authID, newEmail, nil
}

func (c *tokenCache) CreatePasswordReset(ctx context.Context, data *models.CreatePasswordResetToken) *ce.Error {
	pares := constants.CachePrefixPasswordReset
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
		ctx, "s:crepar", script,
		[]string{
			pares + ":" + strconv.FormatUint(data.AuthID, 10), // KEYS[1]
			pares + ":" + data.Token,                          // KEYS[2]
			pares,                                             // KEYS[3]
		},
		[]any{
			data.Token,                   // ARGV[1]
			data.AuthID,                  // ARGV[2]
			int(data.Duration.Seconds()), // ARGV[3]
		},
	)
	if err != nil {
		return ce.NewError(
			ce.CodeCacheScriptExec,
			ce.MsgInternalServer,
			fmt.Errorf("failed to create password reset token: %w", err),
		)
	}

	return nil
}

func (c *tokenCache) CreateVerification(ctx context.Context, data *models.CreateVerificationToken) *ce.Error {
	emver := constants.CachePrefixEmailVerification
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
			emver + ":" + strconv.FormatUint(data.AuthID, 10), // KEYS[1]
			emver + ":" + data.Token,                          // KEYS[2]
			emver,                                             // KEYS[3]
		},
		[]any{
			data.Token,                   // ARGV[1]
			data.AuthID,                  // ARGV[2]
			int(data.Duration.Seconds()), // ARGV[3]
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

func (c *tokenCache) UseVerification(ctx context.Context, token string) (uint64, *ce.Error) {
	emver := constants.CachePrefixEmailVerification
	script := `
		local authID = redis.call("GET", KEYS[1])
		if authID then
			redis.call("DEL", KEYS[1])
			redis.call("DEL", KEYS[2] .. ":" .. authID)
			return authID
		end
		return nil
	`

	res, err := c.cache.Evaluate(
		ctx, "s:usever", script,
		[]string{
			emver + ":" + token, // KEYS[1]
			emver,               // KEYS[2]
		},
		nil,
	)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to use verification token: %w", err)
		if errors.Is(err, ce.ErrCacheNoResult) {
			return 0, ce.NewError(
				ce.CodeInvalidToken,
				ce.MsgInvalidToken,
				wrappedErr,
			)
		}
		return 0, ce.NewError(
			ce.CodeCacheScriptExec,
			ce.MsgInternalServer,
			wrappedErr,
		)
	}

	authID, err := utils.ToUint64(res)
	if err != nil {
		return 0, ce.NewError(
			ce.CodeTypeConversionFailed,
			ce.MsgInternalServer,
			fmt.Errorf("failed to use verification token: %w", err),
		)
	}

	return authID, nil
}
