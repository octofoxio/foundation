/*
 * Copyright (c) 2019. Octofox.io
 */

package logger

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewGlobalLogger(t *testing.T) {
	var (
		log                            = newLogger("Test")
		serviceLogger                  = log.WithServiceID("wallet")
		serviceInfoLogger              = serviceLogger.WithServiceInfo("defaultStellarWallet.GetAsset")
		serviceInfoUserLogger          = serviceInfoLogger.WithUserID("N0091822")
		serviceInfoUserWithFieldLogger = serviceInfoUserLogger.WithField("TransactionID", "00193-aie921jfe-1291jr291j21")
		withURL                        = log.WithURL("GET", "/v1/getsomething").WithRequestID("880-122B")
	)
	t.Run("basic", func(t *testing.T) {
		withURL.Printf("HI\n")
		log.Printf("Hi")
		serviceLogger.Println("Hi")
		serviceInfoLogger.Println("Hi")
		serviceInfoUserLogger.Println("Hi ja")
		serviceInfoUserWithFieldLogger.WithField("RemittanceID", "009182JBBAS831").Println("Hi ja kub")
		serviceLogger.WithRequestID("37e28213684b755066033e7abca4ccf3").Println("With request ID")
		serviceInfoUserLogger.WithField("method", "test").Info("test method")
		serviceInfoUserLogger.WithField("req", "test").Info("test req, must not contain method")
		err := errors.New("some error")
		serviceInfoUserLogger.WithError(err).Error("This is some error description, must not contain other fields")
		serviceInfoUserLogger.Info("log info, must not contain error")
	})

	t.Run("It should use the same instance of formatter but different logger", func(t *testing.T) {
		var loggerInstanceA = serviceInfoUserWithFieldLogger.FieldLogger.(*logrus.Entry).Logger.Formatter
		var loggerInstanceB = log.FieldLogger.(*logrus.Entry).Logger.Formatter
		assert.Equal(t, loggerInstanceA, loggerInstanceB)
		assert.NotEqual(t, serviceInfoUserWithFieldLogger.FieldLogger, log.FieldLogger)
		log.WithField("", "")
	})

	t.Run("with fields should work fine", func(t *testing.T) {

		l := newLogger("gg").WithField("a", "b")
		//.WithField("d", "e")
		l.Info(("HI"))

	})
}

func TestConcurrent(t *testing.T) {
	log := newLogger("concurrency")

	routine := func(name string) {
		for i := 0; i < 10; i++ {
			log.WithField(name, i).Infof("log")
		}
	}

	for i := 0; i < 1000; i++ {
		go routine(fmt.Sprintf("t%d", i))
	}

	time.Sleep(2 * time.Second)
}
