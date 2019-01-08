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
	"os"

	"github.com/sirupsen/logrus"
)

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

// SetLogPattern set logger pattern string
func SetLogPattern(pattern string) error {
	formatter, err := NewLayoutFormatter(pattern)
	if err != nil {
		return err
	}

	isCallerEnable = hasCallerField(formatter.c)
	log.Formatter = formatter
	return nil
}

var (
	isCallerEnable = true
)

func init() {
	log.Out = os.Stdout
	log.Level = logrus.InfoLevel
	log.AddHook(callerHook{})

	SetLogPattern("[%d] [%level] [%P:%F:%M:%L] %msg\n")
}
