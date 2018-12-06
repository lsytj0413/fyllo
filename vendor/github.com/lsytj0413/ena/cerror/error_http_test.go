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
	"testing"

	"github.com/stretchr/testify/suite"
)

type errorHttpTestSuite struct {
	suite.Suite
}

var templateStatus = map[int]int{
	100: 100,
	200: 200,
}

func (s *errorHttpTestSuite) SetupTest() {
	errorsStatus = templateStatus
}

func (s *errorHttpTestSuite) TearDownTest() {
	errorsStatus = map[int]int{}
}

func (s *errorHttpTestSuite) TestStatusCodeFind() {
	for k, v := range templateStatus {
		e := NewRequestError(k, "")
		v1 := e.StatusCode()

		s.Equal(v, v1)
	}
}

func (s *errorHttpTestSuite) TestStatusCodeNotFind() {
	e := NewRequestError(9932121, "")
	v := e.StatusCode()

	s.Equal(http.StatusBadRequest, v)
}

func (s *errorHttpTestSuite) TestSetErrorStatusOK() {
	errorsStatus = map[int]int{}
	SetErrorsStatus(templateStatus)

	s.Equal(len(errorsStatus), len(templateStatus))
	for k, v := range templateStatus {
		v1, ok := errorsStatus[k]
		s.True(ok)
		s.Equal(v, v1)
	}
}

func (s *errorTestSuite) TestSetErrorStatusReplace() {
	errorsStatus = map[int]int{}
	SetErrorsStatus(templateStatus)

	otherStatus := map[int]int{
		200:         100,
		EcodeNotDir: 9987,
	}
	SetErrorsStatus(otherStatus)

	s.Equal(len(templateStatus)+1, len(errorsStatus))
	for k, v := range templateStatus {
		v1, ok := errorsStatus[k]
		s.True(ok)

		_, ok = otherStatus[k]
		if !ok {
			s.Equal(v, v1)
		}
	}

	for k, v := range otherStatus {
		v1, ok := errorsStatus[k]
		s.True(ok)
		s.Equal(v, v1)
	}
}

func TestErrorHttpTestSuite(t *testing.T) {
	s := &errorHttpTestSuite{}
	suite.Run(t, s)
}
