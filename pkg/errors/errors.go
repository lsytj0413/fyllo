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

// Package errors describes errors in project.
package errors

import (
	"net/http"

	"github.com/lsytj0413/ena/cerror"
)

const (
	// EcodeRequestParam errors for Request Param error info
	EcodeRequestParam = 10000001
	// EcodeInitFailed errors for system init error
	EcodeInitFailed = 30000001
	// EcodeProviderNotImplement errors for factory doesn't support plugin
	EcodeProviderNotImplement = 40000002

	// EcodeSequenceOutOfRange errors for sequence exhaust
	EcodeSequenceOutOfRange = 40001001
	// EcodeTimestampRewind errors for current timestamp less than last
	EcodeTimestampRewind = 40001002
	// EcodeSegmentRangeFailed errors for segment item failed, for example the min value bigger than max value
	EcodeSegmentRangeFailed = 40002001
	// EcodeSegmentQueryFailed errors for segment storage failed
	EcodeSegmentQueryFailed = 40002002

	// EcodeInternalError errors for internal server error
	EcodeInternalError = 99999998
	// EcodeUnknown errors for unexpected server error
	EcodeUnknown = 99999999
)

var errorsMessage = map[int]string{
	EcodeRequestParam:         "Request Param Error",
	EcodeInitFailed:           "Server Startup Failed",
	EcodeUnknown:              "Server Unknown Error",
	EcodeProviderNotImplement: "Provider Not Implement",
	EcodeSequenceOutOfRange:   "Sequence Out Of Range",
	EcodeTimestampRewind:      "Current Timestamp Less Than Last",
	EcodeSegmentRangeFailed:   "Current Min Value Bigger Than Max Value",
	EcodeSegmentQueryFailed:   "Query Storage Error",
	EcodeInternalError:        "Internal Server Error",
}

var errorsStatus = map[int]int{
	EcodeUnknown: http.StatusInternalServerError,
}

// NewError construct a cerror.Error and return it
func NewError(errorCode int, cause string) *cerror.Error {
	return cerror.NewError(errorCode, cause)
}

// IsError for cerror.IsError
var IsError = cerror.IsError

// Is for cerror.Is
var Is = cerror.Is

// WriteTo write error message to http response
func WriteTo(w cerror.Writer, err error) error {
	var cerr *cerror.Error
	if IsError(err) {
		cerr = err.(*cerror.Error)
	} else {
		cerr = NewError(EcodeUnknown, err.Error())
	}

	return cerr.WriteTo(w)
}

func init() {
	cerror.SetErrorsMessage(errorsMessage)
	cerror.SetErrorsStatus(errorsStatus)
}
