/*
 * Copyright (c) 2019. Octofox.io
 */

package logger

import (
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
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
	withURL.Printf("HI\n")
	log.Printf("Hi")
	serviceLogger.Println("Hi")
	serviceInfoLogger.Println("Hi")
	serviceInfoUserLogger.Println("Hi ja")
	serviceInfoUserWithFieldLogger.WithField("RemittanceID", "009182JBBAS831").Println("Hi ja kub")
	serviceLogger.WithRequestID("37e28213684b755066033e7abca4ccf3").Println("With request ID")
	err := errors.New("Some error")
	serviceInfoUserLogger.WithError(err).Error("This is some error description")

	t.Run("It should use the same instance of formatter but different logger", func(t *testing.T) {
		var loggerInstanceA = serviceInfoUserWithFieldLogger.FieldLogger.(*logrus.Entry).Logger.Formatter
		var loggerInstanceB = log.FieldLogger.(*logrus.Logger).Formatter
		assert.Equal(t, loggerInstanceA, loggerInstanceB)
		assert.NotEqual(t, serviceInfoUserWithFieldLogger.FieldLogger, log.FieldLogger)
		log.WithField("", "")
	})
}
