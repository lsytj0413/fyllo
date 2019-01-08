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

// FieldName for convert field full name, human readable
type FieldName string

const (
	// FieldNameDate for date field name
	FieldNameDate FieldName = "date"
	// FieldNameLevel for level field name
	FieldNameLevel FieldName = "level"
	// FieldNamePackage for pacakge field name
	FieldNamePackage FieldName = "package"
	// FieldNameMethod for method field name
	FieldNameMethod FieldName = "method"
	// FieldNameFile for file field name
	FieldNameFile FieldName = "file"
	// FieldNameLine for line field name
	FieldNameLine FieldName = "line"
	// FieldNameMessage for message field name
	FieldNameMessage FieldName = "message"

	// FieldNameText for text field name
	FieldNameText FieldName = "text"
)

// FieldKey for convert field key, in layout string
type FieldKey string

const (
	// FieldKeyDate for date field key
	FieldKeyDate FieldKey = "d"
	// FieldKeyLevel for level field key
	FieldKeyLevel FieldKey = "level"
	// FieldKeyPackage for package field key
	FieldKeyPackage FieldKey = "P"
	// FieldKeyMethod for method field key
	FieldKeyMethod FieldKey = "M"
	// FieldKeyFile for file field key
	FieldKeyFile FieldKey = "F"
	// FieldKeyLine for line field key
	FieldKeyLine FieldKey = "L"
	// FieldKeyMessage for message field key
	FieldKeyMessage FieldKey = "msg"

	// FieldKeyText for text field key
	FieldKeyText FieldKey = "text"
)
