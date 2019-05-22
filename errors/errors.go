/*
 * Copyright (c) 2019. Inception Asia
 * Maintain by DigithunWorldwide ❤
 * Maintainer
 * - rungsikorn.r@digithunworldwide.com
 * - nipon.chi@digithunworldwide.com
 * - mai@digithunworldwide.com
 */

package errors

import (
	"encoding/json"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ErrorType int

const (
	ErrorTypeAuth ErrorType = iota
	ErrorTypeBadInput
	ErrorTypeInternal
	ErrorTypeForbidden
	ErrorTypeNotfound
)

type SolarErr struct {
	errorsType ErrorType
	message    string
	detail     []string
	debug      interface{}
	stack      interface{}
}

func (g *SolarErr) Type() ErrorType {
	return g.errorsType
}

func (g *SolarErr) Error() string {
	return g.message
}

func (g *SolarErr) GetDetail() []string {
	return g.detail
}

func (g *SolarErr) GetDebug() interface{} {
	if g.debug == nil {
		return g.stack
	}
	return g.debug
}

func (g SolarErr) WithDetail(detail string) *SolarErr {
	g.detail = append(g.detail, detail)
	return &g
}

func (g SolarErr) WithDebug(d interface{}) *SolarErr {
	g.debug = d
	return &g
}

func NewGlobErr(t ErrorType, message string) *SolarErr {
	return &SolarErr{
		errorsType: t,
		message:    message,
	}
}

func ToGRPCError(err *SolarErr) *status.Status {
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

func FromGRPCError(err error) (*SolarErr, bool) {
	if st, ok := status.FromError(err); ok {
		var data struct {
			ErrorType ErrorType
			Message   string
			Detail    []string
			Debug     interface{}
		}
		if err := json.Unmarshal([]byte(st.Message()), &data); err != nil {
			return nil, false
		}
		return &SolarErr{
			errorsType: data.ErrorType,
			detail:     data.Detail,
			debug:      data.Debug,
			message:    data.Message,
		}, true
	} else {
		return nil, false
	}
}
