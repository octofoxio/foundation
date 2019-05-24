/*
 * Copyright (c) 2019. Octofox.io
 */

package logger

import (
	"fmt"
	"testing"
)

type XX struct {
	*Logger
}

func TestX(t *testing.T) {
	l := New("solarLogger")
	l.Debug("debug message")
	l.Info("info message")
	l.Warn("warn message")
	l.Print("printToBuffer message")
	l.Error("error message")
	l.WithError(fmt.Errorf("something wrongg")).Error("error message again")
	l.WithField("key xx", "value XX").Info("info message")

	WithField("key xx", "value XX").Info("info message")
	WithError(fmt.Errorf("something wrongg")).Error("error CC message")
	Debug("debug message")
	Error("error message")

	xx := XX{
		newLogger("XX"),
	}
	xx.Warn("warning from xx")
	xx.WithError(fmt.Errorf("something wrongg")).Error("error message xx")
	xx.Name = "XX"
	xx.Warn("warning again from xx")
	xx.WithError(fmt.Errorf("something wrongg")).Error("error message xx again")
	xx.WithField("key xx", "value XX").Info("info message")

}
