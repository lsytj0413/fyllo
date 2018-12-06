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

// LogLevel type define, alias for uint8
type LogLevel uint8

const (
	// DebugLevel logs
	DebugLevel LogLevel = iota

	// InfoLevel logs
	InfoLevel

	// WarnLevel logs
	WarnLevel

	// ErrorLevel logs
	ErrorLevel

	// CriticalLevel logs
	CriticalLevel
)
