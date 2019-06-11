/*
 * Copyright (c) 2019. Octofox.io
 */

package foundation

import (
	"context"
	"github.com/octofoxio/foundation/logger"
	"github.com/rs/xid"
)

const (
	FoundationAccessTokenContextKey = "accesstoken"
	FoundationRequestIDContextKey   = "requestid"
	FoundationLoggerContextKey      = "logger"
	FoundationUserIdContextKey      = "userid"
	FoundationMethodContextKey      = "method"
)

func AppendMethodToContext(ctx context.Context, method string, path string) context.Context {
	ctx = context.WithValue(ctx, FoundationMethodContextKey, method)
	var log = GetLoggerFromContext(ctx)
	log = log.WithURL(method, path)
	ctx = AppendLoggerToContext(ctx, log)
	return ctx
}
func AppendUserIDToContext(ctx context.Context, userID string) context.Context {
	ctx = context.WithValue(ctx, FoundationUserIdContextKey, userID)

	var log = GetLoggerFromContext(ctx)
	log = log.WithUserID(userID)

	ctx = AppendLoggerToContext(ctx, log)
	return ctx
}

func GetUserIDFromContext(ctx context.Context) string {
	if userID, ok := ctx.Value(FoundationUserIdContextKey).(string); ok {
		return userID
	} else {
		return ""
	}
}

func AppendLoggerToContext(ctx context.Context, log *logger.Logger) context.Context {
	return context.WithValue(ctx, FoundationLoggerContextKey, log)
}
func GetLoggerFromContext(ctx context.Context) *logger.Logger {
	if log, ok := ctx.Value(FoundationLoggerContextKey).(*logger.Logger); ok && log != nil {
		return log
	} else {
		log = logger.New("foundation")
		log.Warn("You are get logger from context but it empty, make sure you are using foundation.context and append logger before retrieve it")
		return log
	}
}

func GetAccessTokenFromContext(ctx context.Context) string {
	if token, ok := ctx.Value(FoundationAccessTokenContextKey).(string); ok {
		return token
	} else {
		return ""
	}
}
func GetRequestIDFromContext(ctx context.Context) string {
	if token, ok := ctx.Value(FoundationRequestIDContextKey).(string); ok {
		return token
	} else {
		return ""
	}
}

func NewContext(ctx context.Context) context.Context {
	if requestID, ok := ctx.Value(FoundationRequestIDContextKey).(string); !ok || requestID == "" {
		ctx = context.WithValue(ctx, FoundationRequestIDContextKey, xid.New().String())
	}
	var log = logger.New("foundation").WithRequestID(GetRequestIDFromContext(ctx))
	ctx = AppendLoggerToContext(ctx, log)
	return ctx
}
