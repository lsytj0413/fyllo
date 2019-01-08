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

package logger

import (
	"github.com/sirupsen/logrus"

	"github.com/lsytj0413/ena/logger/convert"
)

// LayoutFormatter ...
// %d dateformat
// %level
// %M method name
// %L line
// %msg
// %P package name
// %F file name
type LayoutFormatter struct {
	pattern string
	c       convert.Converter
}

// hasCallerField chech the format pattern is have caller info
func hasCallerField(c convert.Converter) bool {
	hasField := false
	v := convert.NewFieldVisitorForFunc(func(field convert.FieldConverter) bool {
		if field.Properties().IsCallerField {
			hasField = true
			return false
		}
		return true
	})
	c.Visit(v)

	return hasField
}

// NewLayoutFormatter ...
func NewLayoutFormatter(pattern string) (*LayoutFormatter, error) {
	builder := convert.NewBuilder()
	converter, err := builder.Build(pattern)
	if err != nil {
		return nil, err
	}

	return &LayoutFormatter{
		pattern: pattern,
		c:       converter,
	}, nil
}

// Format ...
func (l *LayoutFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	innerEntry := &convert.Entry{
		Time:    entry.Time,
		Message: entry.Message,
		Level:   entry.Level.String(),
		Package: placeHolder,
		File:    placeHolder,
		Method:  placeHolder,
		Line:    placeHolder,
	}

	if d, ok := entry.Data[loggerCallerKeyName]; ok {
		s := d.(*source)
		innerEntry.Package = s.p
		innerEntry.File = s.f
		innerEntry.Method = s.m
		innerEntry.Line = s.l
	}

	return []byte(l.c.Convert(innerEntry)), nil
}
