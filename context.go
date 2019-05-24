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
	FOUNDATION_ACCESS_TOKEN_CONTEXT_KEY = "accesstoken"
	FOUNDATION_REQUEST_ID_CONTEXT_KEY   = "requestid"
	FOUNDATION_LOGGER_CONTEXT_KEY       = "logger"
	FOUNDATION_USER_ID_CONTEXT_KEY      = "userid"
)

func AppendUserIDToContext(ctx context.Context, userID string) context.Context {
	ctx = context.WithValue(ctx, FOUNDATION_USER_ID_CONTEXT_KEY, userID)
	var log = GetLoggerFromContext(ctx)
	log = log.WithUserID(userID)
	ctx = AppendLoggerToContext(ctx, log)
	return ctx
}

func GetUserIDFromContext(ctx context.Context) string {
	if userID, ok := ctx.Value(FOUNDATION_USER_ID_CONTEXT_KEY).(string); ok {
		return userID
	} else {
		return ""
	}
}

func AppendLoggerToContext(ctx context.Context, log *logger.Logger) context.Context {
	return context.WithValue(ctx, FOUNDATION_LOGGER_CONTEXT_KEY, log)
}
func GetLoggerFromContext(ctx context.Context) *logger.Logger {
	if log, ok := ctx.Value(FOUNDATION_LOGGER_CONTEXT_KEY).(*logger.Logger); ok && log != nil {
		return log
	} else {
		log = logger.New("foundation")
		log.Warn("You are get logger from context but it empty, make sure you are using foundation.context and append logger before retrieve it")
		return log
	}
}

func GetAccessTokenFromContext(ctx context.Context) string {
	if token, ok := ctx.Value(FOUNDATION_ACCESS_TOKEN_CONTEXT_KEY).(string); ok {
		return token
	} else {
		return ""
	}
}
func GetRequestIDFromContext(ctx context.Context) string {
	if token, ok := ctx.Value(FOUNDATION_REQUEST_ID_CONTEXT_KEY).(string); ok {
		return token
	} else {
		return ""
	}
}

func NewContext(ctx context.Context) context.Context {
	if requestID, ok := ctx.Value(FOUNDATION_REQUEST_ID_CONTEXT_KEY).(string); !ok || requestID == "" {
		ctx = context.WithValue(ctx, FOUNDATION_REQUEST_ID_CONTEXT_KEY, xid.New().String())
	}
	var log = logger.New("foundation").WithRequestID(GetRequestIDFromContext(ctx))
	ctx = AppendLoggerToContext(ctx, log)
	return ctx
}
