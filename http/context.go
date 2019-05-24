/*
 * Copyright (c) 2019. Octofox.io
 */

package http

import (
	"context"
	"net/http"
)

func getRequestFromContext(c context.Context) *http.Request {
	if req, ok := c.Value(RequestContextKey).(*http.Request); ok {
		return req
	}
	panic("Cannot get request object from context, this is fatal error please ensure you are doing right")
}

func getResponseWriterFromContext(c context.Context) http.ResponseWriter {
	if w, ok := c.Value(ResponseWriterContextKey).(http.ResponseWriter); ok {
		return w
	}
	panic("Cannot get response writer object from context, this is fatal error please ensure you are doing right")
}
