package processor

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/ritchieridanko/apotekly/services/notification/internal/dispatcher"
	"github.com/ritchieridanko/apotekly/services/shared/infra/logger"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
)

const (
	idleDelay  time.Duration = time.Second
	errorDelay time.Duration = 2 * time.Second
)

type Processor struct {
	name       string
	dispatcher *dispatcher.Dispatcher
	logger     *logger.Logger
}

func Init(name string, d *dispatcher.Dispatcher, l *logger.Logger) *Processor {
	return &Processor{name: name, dispatcher: d, logger: l}
}

func (p *Processor) Run(ctx context.Context) error {
	for {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("failed to continue processing pending events (processor=%s): %w", p.name, err)
		}

		err := p.dispatch(ctx)

		switch {
		case err == nil:
			continue
		case err.Code() == ce.CodeNoPendingEventInbox:
			if !p.sleep(ctx, idleDelay) {
				return fmt.Errorf("failed to continue processing pending events (processor=%s): %w", p.name, ctx.Err())
			}
		default:
			p.logger.Error(
				ctx,
				"PENDING EVENT PROCESSING ERROR",
				append(
					err.Fields(),
					logger.NewField("processor", p.name),
					logger.NewField("error_code", err.Code()),
					logger.NewField("error", err.Error()),
				)...,
			)
			if !p.sleep(ctx, errorDelay) {
				return fmt.Errorf("failed to continue processing pending events (processor=%s): %w", p.name, ctx.Err())
			}
		}
	}
}

func (p *Processor) dispatch(ctx context.Context) (err *ce.Error) {
	defer func() {
		err = p.recover(ctx)
	}()

	return p.dispatcher.Dispatch(ctx)
}

func (p *Processor) recover(ctx context.Context) *ce.Error {
	if r := recover(); r != nil {
		p.logger.Error(
			ctx,
			"PANIC RECOVERED",
			logger.NewField("processor", p.name),
			logger.NewField("panic", fmt.Sprintf("%v", r)),
			logger.NewField("stack_trace", string(debug.Stack())),
		)
		return ce.NewError(
			ce.CodePanicOccurred,
			ce.MsgInternalServer,
			fmt.Errorf("%v", r),
		)
	}
	return nil
}

func (p *Processor) sleep(ctx context.Context, d time.Duration) bool {
	timer := time.NewTimer(d)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		return true
	}
}
