package cache

import (
	"context"
	"fmt"
	"strconv"

	"github.com/ritchieridanko/apotekly/services/gateway/internal/constants"
	"github.com/ritchieridanko/apotekly/services/gateway/internal/models"
	cc "github.com/ritchieridanko/apotekly/services/shared/infra/cache"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
)

type RateLimiterCache interface {
	LimitByFixedWindow(ctx context.Context, data *models.LimitByFixedWindowRL) (err *ce.Error)
	LimitBySlidingWindow(ctx context.Context, data *models.LimitBySlidingWindowRL) (err *ce.Error)
}

type rateLimiterCache struct {
	cache *cc.Cache
}

func NewRateLimiterCache(cc *cc.Cache) RateLimiterCache {
	return &rateLimiterCache{cache: cc}
}

func (c *rateLimiterCache) LimitByFixedWindow(ctx context.Context, data *models.LimitByFixedWindowRL) *ce.Error {
	script := `
		local current = redis.call("GET", KEYS[1])
		if current and tonumber(current) >= tonumber(ARGV[2]) then
			return 0
		end

		local tracker = redis.call("INCR", KEYS[1])
		if tonumber(tracker) == 1 then
			redis.call("EXPIRE", KEYS[1], ARGV[1])
		end

		return 1
	`

	res, err := c.cache.Evaluate(
		ctx, "s:lbfw", script,
		[]string{
			constants.CachePrefixRateLimit + ":" + data.Namespace + ":" + data.IPAddress, // KEYS[1]
		},
		[]any{
			int(data.Window.Seconds()), // ARGV[1]
			data.Limit,                 // ARGV[2]
		},
	)
	if err != nil {
		return ce.NewError(
			ce.CodeCacheScriptExec,
			ce.MsgInternalServer,
			fmt.Errorf("failed to rate limit by fixed window: %w", err),
		)
	}
	if res.(int64) == 0 {
		return ce.NewError(ce.CodeTooManyRequests, ce.MsgTooManyRequests, nil)
	}

	return nil
}

func (c *rateLimiterCache) LimitBySlidingWindow(ctx context.Context, data *models.LimitBySlidingWindowRL) *ce.Error {
	script := `
		local now = tonumber(ARGV[1])
		local window = tonumber(ARGV[2])

		local curr_window = math.floor(now / window)
		local curr_key = KEYS[1] .. ":" .. curr_window
		local prev_key = KEYS[1] .. ":" .. (curr_window - 1)

		local curr_count = tonumber(redis.call("GET", curr_key) or 0)
		local prev_count = tonumber(redis.call("GET", prev_key) or 0)

		-- Calculate weight of the previous window based on elapsed time
		local time_into_window = now % window
		local weight = (window - time_into_window) / window
		local estimated_count = math.floor(prev_count * weight + curr_count)

		if estimated_count >= tonumber(ARGV[3]) then
			return 0
		end

		local tracker = redis.call("INCR", curr_key)
		if tonumber(tracker) == 1 then
			redis.call("EXPIRE", curr_key, window * 2)
		end

		return 1
	`

	res, err := c.cache.Evaluate(
		ctx, "s:lbsw", script,
		[]string{
			constants.CachePrefixRateLimit + ":" + data.Namespace + ":" + strconv.FormatUint(data.AuthID, 10), // KEYS[1]
		},
		[]any{
			data.Now.Unix(),            // ARGV[1]
			int(data.Window.Seconds()), // ARGV[2]
			data.Limit,                 // ARGV[3]
		},
	)
	if err != nil {
		return ce.NewError(
			ce.CodeCacheScriptExec,
			ce.MsgInternalServer,
			fmt.Errorf("failed to rate limit by sliding window: %w", err),
		)
	}
	if res.(int64) == 0 {
		return ce.NewError(ce.CodeTooManyRequests, ce.MsgTooManyRequests, nil)
	}

	return nil
}
