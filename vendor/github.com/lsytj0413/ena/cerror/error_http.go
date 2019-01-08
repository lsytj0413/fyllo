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

package cerror

import (
	"net/http"
)

var errorsStatus = map[int]int{}

// NewRequestError construct a Request Error struct
func NewRequestError(errorCode int, cause string) *Error {
	return NewError(errorCode, cause)
}

// StatusCode returns the RequestError.httpStatusCode
func (e Error) StatusCode() int {
	status, ok := errorsStatus[e.ErrorCode]
	if !ok {
		return http.StatusBadRequest
	}

	return status
}

// Writer is a interface define for write to http.Response
type Writer interface {
	WriteHeader(int)
	Write([]byte) (int, error)
}

// WriteTo write error message to http response
func (e Error) WriteTo(w Writer) error {
	w.WriteHeader(e.StatusCode())
	_, err := w.Write([]byte(e.JSONString() + "\n"))
	return err
}

// SetErrorsStatus init error defined errorCode and httpStatusCode
func SetErrorsStatus(status map[int]int) {
	for k, v := range status {
		errorsStatus[k] = v
	}
}
