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

// Package cerror defines a common error information defines
package cerror

import (
	"encoding/json"
	"fmt"
)

// Error is store package error message define
type Error struct {
	ErrorCode int    `json:"errorCode"`
	Message   string `json:"message"`
	Cause     string `json:"cause,omitempty"`
}

var errorsMessage = map[int]string{}

// NewError construct a Error struct and return it
func NewError(errorCode int, cause string) *Error {
	return &Error{
		ErrorCode: errorCode,
		Message:   errorsMessage[errorCode],
		Cause:     cause,
	}
}

// Error is for the error interface
func (e Error) Error() string {
	return e.Message + " (" + e.Cause + ")"
}

var (
	// For unittest
	marshal func(interface{}) ([]byte, error)
)

// JSONString returns the JSON format message
func (e Error) JSONString() string {
	b, err := marshal(e)
	if err != nil {
		return fmt.Sprintf(
			`{"errorCode":%d,"message":"%s","cause":"%s"}`,
			e.ErrorCode,
			e.Message,
			e.Cause)
	}

	return string(b)
}

// SetErrorsMessage init error defined errorCode and Message
func SetErrorsMessage(message map[int]string) {
	for k, v := range message {
		errorsMessage[k] = v
	}
}

func init() {
	marshal = json.Marshal
}
