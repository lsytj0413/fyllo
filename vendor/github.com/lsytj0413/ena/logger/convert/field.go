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

// FieldProperties for convert field properties
type FieldProperties struct {
	IsCallerField bool
}

// FieldConverter convert field to string
type FieldConverter interface {
	Name() string
	Key() string
	Convert(entry *Entry) string
	Properties() FieldProperties
}

type callerFieldProperties struct {
}

func (c *callerFieldProperties) Properties() FieldProperties {
	return FieldProperties{
		IsCallerField: true,
	}
}

type nonCallerFieldProperties struct {
}

func (c *nonCallerFieldProperties) Properties() FieldProperties {
	return FieldProperties{
		IsCallerField: false,
	}
}

// textConvert for const text
type textConverter struct {
	nonCallerFieldProperties
	text string
}

func (c *textConverter) Name() string {
	return string(FieldNameText)
}

func (c *textConverter) Key() string {
	return string(FieldKeyText)
}

func (c *textConverter) Convert(entry *Entry) string {
	return c.text
}

// dateConverter for date field
type dateConverter struct {
	nonCallerFieldProperties
	timestampFormat string
}

func (c *dateConverter) Name() string {
	return string(FieldNameDate)
}

func (c *dateConverter) Key() string {
	return string(FieldKeyDate)
}

func (c *dateConverter) Convert(entry *Entry) string {
	return entry.Time.Format(c.timestampFormat)
}

// levelConverter for level field
type levelConverter struct {
	nonCallerFieldProperties
}

func (c *levelConverter) Name() string {
	return string(FieldNameLevel)
}

func (c *levelConverter) Key() string {
	return string(FieldKeyLevel)
}

func (c *levelConverter) Convert(entry *Entry) string {
	return entry.Level
}

// packageConverter for package field
type packageConverter struct {
	callerFieldProperties
}

func (c *packageConverter) Name() string {
	return string(FieldNamePackage)
}

func (c *packageConverter) Key() string {
	return string(FieldKeyPackage)
}

func (c *packageConverter) Convert(entry *Entry) string {
	return entry.Package
}

// fileConverter for file field
type fileConverter struct {
	callerFieldProperties
}

func (c *fileConverter) Name() string {
	return string(FieldNameFile)
}

func (c *fileConverter) Key() string {
	return string(FieldKeyFile)
}

func (c *fileConverter) Convert(entry *Entry) string {
	return entry.File
}

// methodConverter for method field
type methodConverter struct {
	callerFieldProperties
}

func (c *methodConverter) Name() string {
	return string(FieldNameMethod)
}

func (c *methodConverter) Key() string {
	return string(FieldKeyMethod)
}

func (c *methodConverter) Convert(entry *Entry) string {
	return entry.Method
}

// lineConverter for line field
type lineConverter struct {
	callerFieldProperties
}

func (c *lineConverter) Name() string {
	return string(FieldNameLine)
}

func (c *lineConverter) Key() string {
	return string(FieldKeyLine)
}

func (c *lineConverter) Convert(entry *Entry) string {
	return entry.Line
}

// messageConverter for message field
type messageConverter struct {
	nonCallerFieldProperties
}

func (c *messageConverter) Name() string {
	return string(FieldNameMessage)
}

func (c *messageConverter) Key() string {
	return string(FieldKeyMessage)
}

func (c *messageConverter) Convert(entry *Entry) string {
	return entry.Message
}
