/*
 * Copyright (c) 2019. Octofox.io
 */

package app

import (
	"context"
	"fmt"
)

type StringSvc struct{}

func NewStringSvc() *StringSvc {
	return &StringSvc{}
}

func (s *StringSvc) Concat(c context.Context, input *ConcatInput) (*ConcatOutput, error) {
	return &ConcatOutput{
		Result: fmt.Sprintf("%s%s", input.Origin, input.Extend),
	}, nil

}
