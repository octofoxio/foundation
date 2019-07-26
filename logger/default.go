/*
 * Copyright (c) 2019. Octofox.io
 */

package logger

/*
 * Copyright (c) 2019. Inception Asia
 * Maintain by DigithunWorldwide ❤
 * Maintainer
 * - rungsikorn.r@digithunworldwide.com
 * - nipon.chi@digithunworldwide.com
 * - mai@digithunworldwide.com
 */
import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	fieldServiceID   = "service-id"
	fieldServiceInfo = "service-info"
	fieldUserID      = "user-id"
	fieldData        = "data"
	fieldRequestID   = "request-id"
	fieldError       = "error"
	fieldURL         = "url"
)

type globalLogFormatter struct{}

func valueToString(value interface{}) string {
	var v interface{} = nil
	switch value := value.(type) {
	case string:
		v = value
	case error:
		v = value.Error()
	case fmt.Stringer:
		v = value.String()
	case fmt.GoStringer:
		v = value.GoString()
	default:
		v = value
	}
	return fmt.Sprintf("%v", v)
}
func printToBuffer(b *bytes.Buffer, value interface{}, defaultValue string) {
	if value != nil {
		b.WriteString(fmt.Sprintf("| %s ", valueToString(value)))
	} else {
		b.WriteString(fmt.Sprintf("| %s ", defaultValue))
	}
}

// Print log from custom field entry
// <RFC date> | <Level> | <Request ID> <URL> <Name> | <Service ID> | <Library information> "<Message>" "...<Field vey>=<Field value>"
func (g *globalLogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	// format data string
	var Data interface{}
	var toStringData = make([]string, 0, len(entry.Data))
	for k, v := range entry.Data {
		if k == fieldData || k == fieldRequestID || k == fieldServiceID || k == fieldServiceInfo || k == fieldURL || k == fieldUserID {
			continue
		}
		toStringData = append(toStringData, fmt.Sprintf("%s=%s", k, valueToString(v)))
	}
	// print data only with value
	if len(toStringData) > 0 {
		Data = strings.Join(toStringData[:], " | ")
	}

	var (
		RFCDate     = time.Now().Format(time.RFC1123Z + " ")
		ServiceID   = entry.Data[fieldServiceID]
		ServiceInfo = entry.Data[fieldServiceInfo]
		UserID      = entry.Data[fieldUserID]
		RequestID   = entry.Data[fieldRequestID]
		RequestURL  = entry.Data[fieldURL]
	)

	b.WriteString(RFCDate)
	if RequestID != nil {
		b.WriteString(fmt.Sprintf("(%s)", RequestID))
	}
	printToBuffer(b, strings.ToUpper(entry.Level.String()), "INFO")
	printToBuffer(b, UserID, "system")
	if RequestURL != nil && RequestURL != "" {
		b.WriteString(" ")
		b.WriteString(fmt.Sprintf("%s ", RequestURL))
	}
	printToBuffer(b, ServiceID, "-")
	printToBuffer(b, ServiceInfo, "-")
	printToBuffer(b, entry.Message, "-")
	if Data != nil {
		printToBuffer(b, Data, "")
	}
	b.WriteString("\n")
	return b.Bytes(), nil
}
func newPrivateLogger(output io.Writer) *logrus.Logger {
	var log = &logrus.Logger{
		Out:       output,
		Formatter: &globalLogFormatter{},
		Level:     logrus.DebugLevel,
	}
	return log
}
func newLogger(name string) *Logger {
	var log = &Logger{
		FieldLogger: newPrivateLogger(os.Stdout),
		Name:        name,
		Data:        map[string]interface{}{},
		isInitial:   true,
		mux:         &sync.Mutex{},
	}
	return log.WithServiceID(name)
}

type Logger struct {
	Name               string
	ServiceID          string
	ServiceInfo        string
	UserID             string
	RequestID          string
	logrus.FieldLogger // Logger instance
	Data               map[string]interface{}
	isInitial          bool
	mux                *sync.Mutex
}

func (g Logger) SetOutput(w io.Writer) *Logger {
	g.FieldLogger = newPrivateLogger(w).WithFields(g.Data)
	return &g
}

func (g Logger) Printf(format string, args ...interface{}) {
	if !g.isInitial {
		g = *New(g.Name)
	}
	g.FieldLogger.Printf(format, args...)
}

func (g Logger) WithError(err error) *Logger {
	return g.WithField(fieldError, err)
}

func (g Logger) WithField(key string, value interface{}) *Logger {
	if !g.isInitial {
		g = *New(g.Name)
	}
	f := *g.FieldLogger.WithField(key, value)
	if g.mux == nil {
		g.mux = &sync.Mutex{}
	}
	g.mux.Lock()
	g.Data[key] = value
	g.mux.Unlock()
	g.FieldLogger = &f
	return &g
}

func (g Logger) WithRequestID(value string) *Logger {
	return g.WithField(fieldRequestID, value)
}
func (g Logger) WithServiceID(value string) *Logger {
	return g.WithField(fieldServiceID, value)
}

func (g Logger) WithServiceInfo(value string) *Logger {
	return g.WithField(fieldServiceInfo, value)
}

func (g Logger) WithURL(method string, url string) *Logger {
	return g.WithField(fieldURL, fmt.Sprintf("%s %s", method, url))
}

func (g Logger) WithUserID(ID string) *Logger {
	return g.WithField(fieldUserID, ID)
}

func New(name string) *Logger {
	return newLogger(name)
}
