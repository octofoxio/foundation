/*
 * Copyright (c) 2019. Octofox.io
 */

package errors

import (
	"encoding/json"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ErrorType int

const (
	ErrorTypeAuth      ErrorType = 401
	ErrorTypeBadInput            = 400
	ErrorTypeInternal            = 500
	ErrorTypeForbidden           = 403
	ErrorTypeNotfound            = 404
)

type Error struct {
	errorsType ErrorType
	message    string
	detail     []string
	debug      interface{}
	stack      interface{}
}

func (g *Error) Type() ErrorType {
	return g.errorsType
}

func (g *Error) Error() string {
	return g.message
}

func (g *Error) GetDetail() []string {
	return g.detail
}

func (g *Error) GetDebug() interface{} {
	if g.debug == nil {
		return g.stack
	}
	return g.debug
}

func (g Error) WithDetail(detail string) *Error {
	g.detail = append(g.detail, detail)
	return &g
}

func (g Error) WithDebug(d interface{}) *Error {
	g.debug = d
	return &g
}

func NewGlobErr(t ErrorType, message string) *Error {
	return &Error{
		errorsType: t,
		message:    message,
	}
}

func ToGRPCError(err *Error) *status.Status {
	type data struct {
		ErrorType ErrorType   `json:"type"`
		Message   string      `json:"message"`
		Detail    []string    `json:"detail,omitempty"`
		Debug     interface{} `json:"debug,omitempty"`
	}
	b, _ := json.Marshal(&data{
		Message:   err.message,
		Debug:     err.debug,
		Detail:    err.detail,
		ErrorType: err.errorsType,
	})

	// handle Glob Error แล้วเปลี่ยนเป็น GRPC error
	// โดยจะต้องทำให้ error สามารถแปลงกลับไปเป็น GlobError ได้
	// ด้วยการจะแนบ JSON ไปใน message
	switch err.Type() {
	case ErrorTypeBadInput:
		return status.New(codes.InvalidArgument, string(b))
	case ErrorTypeForbidden:
	case ErrorTypeAuth:
		return status.New(codes.PermissionDenied, string(b))
	}
	return status.New(codes.Internal, string(b))
}

func FromGRPCError(err error) (*Error, bool) {
	if st, ok := status.FromError(err); ok {
		var data struct {
			ErrorType ErrorType   `json:"type"`
			Message   string      `json:"message"`
			Detail    []string    `json:"detail,omitempty"`
			Debug     interface{} `json:"debug,omitempty"`
		}
		if err := json.Unmarshal([]byte(st.Message()), &data); err != nil {
			return nil, false
		}
		return &Error{
			errorsType: data.ErrorType,
			detail:     data.Detail,
			debug:      data.Debug,
			message:    data.Message,
		}, true
	} else {
		return nil, false
	}
}
