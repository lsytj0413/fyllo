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

type errorHTTPTestSuite struct {
	suite.Suite
}

var templateStatus = map[int]int{
	100: 100,
	200: 200,
}

func (s *errorHTTPTestSuite) SetupTest() {
	errorsStatus = templateStatus
}

func (s *errorHTTPTestSuite) TearDownTest() {
	errorsStatus = map[int]int{}
}

func (s *errorHTTPTestSuite) TestStatusCodeFind() {
	for k, v := range templateStatus {
		e := NewRequestError(k, "")
		v1 := e.StatusCode()

		s.Equal(v, v1)
	}
}

func (s *errorHTTPTestSuite) TestStatusCodeNotFind() {
	e := NewRequestError(9932121, "")
	v := e.StatusCode()

	s.Equal(http.StatusBadRequest, v)
}

func (s *errorHTTPTestSuite) TestSetErrorStatusOK() {
	errorsStatus = map[int]int{}
	SetErrorsStatus(templateStatus)

	s.Equal(len(errorsStatus), len(templateStatus))
	for k, v := range templateStatus {
		v1, ok := errorsStatus[k]
		s.True(ok)
		s.Equal(v, v1)
	}
}

type fakeWriter struct {
	header int
	body   []byte
}

func (f *fakeWriter) WriteHeader(v int) {
	f.header = v
}

func (f *fakeWriter) Write(b []byte) (int, error) {
	f.body = b
	return 0, nil
}
func (s *errorHTTPTestSuite) TestWriteToOk() {
	err := NewError(0, "cause")
	w := &fakeWriter{}

	err.WriteTo(w)

	s.Equal(err.StatusCode(), w.header)
	s.Equal(err.JSONString()+"\n", string(w.body))
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
	s := &errorHTTPTestSuite{}
	suite.Run(t, s)
}
