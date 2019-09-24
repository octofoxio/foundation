/*
 * Copyright (c) 2019. Octofox.io
 */

package foundationerrorv2

import (
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
)

type Error struct {
	Type    codes.Code
	message []string
}

func (e *Error) Error() string {
	return fmt.Sprintf(strings.Join(e.message, "\n"))
}

func (e *Error) AppendMessage(msg string, args ...interface{}) *Error {
	e.message = append(e.message, fmt.Sprintf(msg, args...))
	return e
}

func FromStatusError(status *status.Status) *Error {
	return &Error{
		Type:    status.Code(),
		message: strings.Split(status.Message(), "\n"),
	}
}

func New(Type codes.Code) *Error {
	return &Error{
		Type: Type,
	}
}
