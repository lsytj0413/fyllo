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

package cmap

import "fmt"

// InvalidParamError is present the parameter illegal
type InvalidParamError struct {
	msg string
}

func newInvalidParamError(msg string) *InvalidParamError {
	return &InvalidParamError{
		msg: fmt.Sprintf("cmap: invalid param(%s)", msg),
	}
}

func (e InvalidParamError) Error() string {
	return e.msg
}

// InvalidPairTypeError is present the invalid Pair key-value type
type InvalidPairTypeError struct {
	msg string
}

func newInvalidPairTypeError(msg string) *InvalidPairTypeError {
	return &InvalidPairTypeError{
		msg: fmt.Sprintf("cmap: invalid pair type(%s)", msg),
	}
}

func (e InvalidPairTypeError) Error() string {
	return e.msg
}

// PairRedistributorError is present the error Pair Redistributor
type PairRedistributorError struct {
	msg string
}

func newPairRedistributorError(msg string) *PairRedistributorError {
	return &PairRedistributorError{
		msg: fmt.Sprintf("cmap: failing pair redistributor(%s)", msg),
	}
}

func (e PairRedistributorError) Error() string {
	return e.msg
}
