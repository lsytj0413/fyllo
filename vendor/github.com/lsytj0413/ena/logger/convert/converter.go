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

package convert

import (
	"strings"
	"time"
)

// FieldVisitor for visit field converters
type FieldVisitor interface {
	// Visit the field converter, and return true for next, return false to stop this time visit
	Visit(FieldConverter) bool
}

type funcFieldVisitor struct {
	fn func(FieldConverter) bool
}

func (v *funcFieldVisitor) Visit(field FieldConverter) bool {
	return v.fn(field)
}

// NewFieldVisitorForFunc return FieldVisitor wrap for function
func NewFieldVisitorForFunc(fn func(FieldConverter) bool) FieldVisitor {
	return &funcFieldVisitor{
		fn: fn,
	}
}

// Converter convert fields to string
type Converter interface {
	Convert(entry *Entry) string
	Visit(v FieldVisitor)
}

// Entry for convert input
type Entry struct {
	Time    time.Time
	Level   string
	Package string
	File    string
	Method  string
	Line    string
	Message string
	Data    map[string]string
}

// defConverterImpl implememnt Convert interface
type defConverterImpl struct {
	fields []FieldConverter
}

func (c *defConverterImpl) Convert(entry *Entry) string {
	buffer := strings.Builder{}
	for _, f := range c.fields {
		buffer.WriteString(f.Convert(entry))
	}
	return buffer.String()
}

func (c *defConverterImpl) Visit(v FieldVisitor) {
	for _, field := range c.fields {
		if !v.Visit(field) {
			return
		}
	}
}

// NewConverter construct Converter
func NewConverter(fields []FieldConverter) Converter {
	return &defConverterImpl{
		fields: fields,
	}
}
