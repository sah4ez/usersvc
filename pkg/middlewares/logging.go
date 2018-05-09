package middlewares

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/sah4ez/usersvc/pkg/service"
)

type Middleware func(service.Service) service.Service

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   service.Service
	logger log.Logger
}

func (mw loggingMiddleware) PostUser(ctx context.Context, p service.User) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PostUser", "id", p.ID, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.PostUser(ctx, p)
}

func (mw loggingMiddleware) GetUser(ctx context.Context, id string) (p service.User, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetUser", "id", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetUser(ctx, id)
}

func (mw loggingMiddleware) PatchUser(ctx context.Context, id string, p service.User) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PatchUser", "id", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.PatchUser(ctx, id, p)
}

func (mw loggingMiddleware) GetUsers(ctx context.Context) (users []service.User, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetUsers", "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetUsers(ctx)
}
