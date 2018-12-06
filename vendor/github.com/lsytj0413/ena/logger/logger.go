// Copyright (c) 2018 soren yang
//
// Licensed under the MIT License
// you may not use this file except in complicance with the License.
// You may obtain a copy of the License at
//
//     https://opensource.org/licenses/MIT
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package logger provide a wrapper of Sirupsen/logrus
package logger

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/Sirupsen/logrus"
)

type contextHook struct {
}

func (contextHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (contextHook) Fire(entry *logrus.Entry) error {
	entry.Data["file"] = "unknown:unknown:0"

	pc := make([]uintptr, 3, 3)
	cnt := runtime.Callers(7, pc)
	found := false
	for i := 0; i < cnt; i++ {
		fu := runtime.FuncForPC(pc[i] - 1)
		name := fu.Name()
		if !strings.Contains(name, "github.com/Sirupsen/logrus") && !found {
			found = true
			continue
		}
		if found {
			file, line := fu.FileLine(pc[i] - 1)
			entry.Data["file"] = fmt.Sprintf("%v:%v:%v", path.Base(file), path.Base(name), line)
			break
		}
	}

	return nil
}

var log = logrus.New()

// Debugf logs DEBUG level to Dest Output
func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

// Infof logs INFO level to Dest Output
func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

// Warnf logs WARN level to Dest Output
func Warnf(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

// Errorf logs ERROR level to Dest Output
func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

// Criticalf logs Critical level to Dest Output
func Criticalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

// Debug logs DEBUG level to Dest Output
func Debug(args ...interface{}) {
	log.Debug(args...)
}

// Info logs INFO level to Dest Output
func Info(args ...interface{}) {
	log.Info(args...)
}

// Warn logs WARN level to Dest Output
func Warn(args ...interface{}) {
	log.Warn(args...)
}

// Error logs ERROR level to Dest Output
func Error(args ...interface{}) {
	log.Error(args...)
}

// Critical logs Critical level to Dest Output
func Critical(args ...interface{}) {
	log.Fatal(args...)
}

// Debugln logs DEBUG level to Dest Output
func Debugln(args ...interface{}) {
	log.Debugln(args...)
}

// Infoln logs INFO level to Dest Output
func Infoln(args ...interface{}) {
	log.Infoln(args...)
}

// Warnln logs WARN level to Dest Output
func Warnln(args ...interface{}) {
	log.Warnln(args...)
}

// Errorln logs ERROR level to Dest Output
func Errorln(args ...interface{}) {
	log.Errorln(args...)
}

// Criticalln logs Critical level to Dest Output
func Criticalln(args ...interface{}) {
	log.Fatalln(args...)
}

// SetLogLevel changes logger.Level, the default logger.Level is InfoLevel
func SetLogLevel(v LogLevel) {
	switch v {
	case DebugLevel:
		log.Level = logrus.DebugLevel
	default:
		fallthrough
	case InfoLevel:
		log.Level = logrus.InfoLevel
	case WarnLevel:
		log.Level = logrus.WarnLevel
	case ErrorLevel:
		log.Level = logrus.ErrorLevel
	case CriticalLevel:
		log.Level = logrus.FatalLevel
	}
}

func init() {
	log.Out = os.Stdout
	log.Level = logrus.InfoLevel
	log.AddHook(contextHook{})

	formatter := logrus.TextFormatter{DisableTimestamp: false,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05.000000000Z07:00"}
	log.Formatter = &formatter
}
