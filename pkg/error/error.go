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

// Package error describes errors in project.
package error

import (
	"net/http"

	"github.com/lsytj0413/ena/cerror"
)

const (
	// EcodeRequestParam errors for Request Param error info
	EcodeRequestParam = 10000001
	// EcodeInitFailed errors for system init error
	EcodeInitFailed = 30000001
	// EcodeNotImplement errors for system not implement
	EcodeNotImplement = 40000001
	// EcodePluginNotImplement errors for factory doesn't support plugin
	EcodePluginNotImplement = 40000002

	// EcodeSequenceOutOfRange errors for sequence exhaust
	EcodeSequenceOutOfRange = 40001001
	// EcodeTimestampRewind errors for current timestamp less than last
	EcodeTimestampRewind = 40001002

	// EcodeUnknown errors for unexpected server error
	EcodeUnknown = 99999999
)

var errorsMessage = map[int]string{
	EcodeRequestParam:       "Request Param Error",
	EcodeInitFailed:         "Server Startup Failed",
	EcodeUnknown:            "Server Unknown Error",
	EcodeNotImplement:       "Not Implement",
	EcodePluginNotImplement: "Plugin Not Implement",
	EcodeSequenceOutOfRange: "Sequence Out Of Range",
	EcodeTimestampRewind:    "Current Timestamp Less Than Last",
}

var errorsStatus = map[int]int{
	EcodeUnknown: http.StatusInternalServerError,
}

// NewError construct a cerror.Error and return it
func NewError(errorCode int, cause string) *cerror.Error {
	return cerror.NewError(errorCode, cause)
}

// IsPluginNotImplement checks the error is EcodePluginNotImplement
func IsPluginNotImplement(err error) bool {
	if ce, ok := err.(*cerror.Error); ok {
		return ce.ErrorCode == EcodePluginNotImplement
	}

	return false
}

// NewPluginNotImplement construct a PluginNotImplementError
func NewPluginNotImplement(cause string) *cerror.Error {
	return NewError(EcodePluginNotImplement, cause)
}

func init() {
	cerror.SetErrorsMessage(errorsMessage)
	cerror.SetErrorsStatus(errorsStatus)
}
