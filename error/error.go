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
	// EcodeUnknown errors for unexpected server error
	EcodeUnknown = 99999999
)

var errorsMessage = map[int]string{
	EcodeRequestParam: "Request Param Error",
	EcodeInitFailed:   "Server Startup Failed",
	EcodeUnknown:      "Server Unknown Error",
}

var errorsStatus = map[int]int{
	EcodeUnknown: http.StatusInternalServerError,
}

// NewError const struct a cerror.Error and return it
func NewError(errorCode int, cause string) *cerror.Error {
	return cerror.NewError(errorCode, cause)
}

func init() {
	cerror.SetErrorsMessage(errorsMessage)
	cerror.SetErrorsStatus(errorsStatus)
}
