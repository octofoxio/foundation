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
	var (
		RFCDate     = time.Now().Format(time.RFC1123Z + " ")
		ServiceID   = entry.Data[fieldServiceID]
		ServiceInfo = entry.Data[fieldServiceInfo]
		UserID      = entry.Data[fieldUserID]
		Data        = entry.Data[fieldData]
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
}

func (g Logger) SetOutput(w io.Writer) *Logger {
	g.FieldLogger = newPrivateLogger(w).WithFields(g.Data)
	return &g
}

func (g Logger) setAttribute(key, value string) *Logger {
	if !g.isInitial {
		g = *New(g.Name)
	}
	f := *g.FieldLogger.WithField(key, value)
	g.Data[key] = value
	g.FieldLogger = &f
	return &g
}

func (g Logger) Printf(format string, args ...interface{}) {
	if !g.isInitial {
		g = *New(g.Name)
	}
	g.FieldLogger.Printf(format, args...)
}

func (g Logger) WithError(err error) *Logger {
	return g.setAttribute(fieldError, value)
}

func (g Logger) WithField(key string, value interface{}) *Logger {
	if !g.isInitial {
		g = *New(g.Name)
	}
	g.Data[key] = value
	var toStringData = make([]string, 0, len(g.Data))
	for k, v := range g.Data {
		if k == fieldData || k == fieldError || k == fieldRequestID || k == fieldServiceID || k == fieldServiceInfo || k == fieldURL || k == fieldUserID {
			continue
		}
		toStringData = append(toStringData, fmt.Sprintf("%s=%s", k, valueToString(v)))
	}
	return g.setAttribute(fieldData, strings.Join(toStringData[:], " | "))
}

func (g Logger) WithRequestID(value string) *Logger {
	return g.setAttribute(fieldRequestID, value)
}
func (g Logger) WithServiceID(value string) *Logger {
	return g.setAttribute(fieldServiceID, value)
}

func (g Logger) WithServiceInfo(value string) *Logger {
	return g.setAttribute(fieldServiceInfo, value)
}

func (g Logger) WithURL(method string, url string) *Logger {
	return g.setAttribute(fieldURL, fmt.Sprintf("%s %s", method, url))
}

func (g Logger) WithUserID(ID string) *Logger {
	return g.setAttribute(fieldUserID, ID)
}

func New(name string) *Logger {
	return newLogger(name)
}
