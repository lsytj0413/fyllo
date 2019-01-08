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

package fs

import (
	"github.com/lsytj0413/ena/cerror"
)

const (
	// EcodeUnknown is unknown error info
	EcodeUnknown = 10009999
	// EcodeNotFile errors for operate on dir but file is required
	EcodeNotFile = 10000001
	// EcodeNotDir errors for operate on file but dir is required
	EcodeNotDir = 10000002
	// EcodeNotExists errors for operate on target but doesn't exists
	EcodeNotExists = 10000003
	// EcodeExists errors for Add target but already exists
	EcodeExists = 10000004
	// EcodeDirNotEmpty errors for Remove directory but directory has child etc
	EcodeDirNotEmpty = 10000005
)

var errorsMessage = map[int]string{
	EcodeUnknown:     "Unknown Error",
	EcodeNotFile:     "Target is Not File",
	EcodeNotDir:      "Target is Not Dir",
	EcodeNotExists:   "Target is not exists",
	EcodeExists:      "Target is exists",
	EcodeDirNotEmpty: "Dir not empty",
}

// NewError construct a Error struct and return it
func NewError(errorCode int, cause string) *cerror.Error {
	return cerror.NewError(errorCode, cause)
}

// Is check the error type and errorCode
var Is = cerror.Is

func init() {
	cerror.SetErrorsMessage(errorsMessage)
}
