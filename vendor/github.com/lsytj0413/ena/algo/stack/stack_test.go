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

package stack

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type stackTestSuite struct {
	suite.Suite

	s Stack
}

func (s *stackTestSuite) SetupTest() {
	s.s = New()
}

func (s *stackTestSuite) TestPushOk() {
	r := s.s.Peek()
	s.Nil(r)

	values := []int{1, 2, 3, 4, 5}
	for _, v := range values {
		s.s.Push(v)
		r = s.s.Peek()
		s.Equal(v, r.(int))
	}
}

func (s *stackTestSuite) TestPopNil() {
	r := s.s.Pop()
	s.Nil(r)
}

func (s *stackTestSuite) TestPopOk() {
	values := []int{1, 2, 3, 4, 5}
	for _, v := range values {
		s.s.Push(v)
	}

	for i := len(values) - 1; i >= 0; i-- {
		v := s.s.Pop()
		s.Equal(values[i], v.(int))
	}
}

func (s *stackTestSuite) TestPeekNil() {
	r := s.s.Peek()
	s.Nil(r)
}

func (s *stackTestSuite) TestPeekOk() {
	values := []int{1, 2, 3, 4, 5}
	for _, v := range values {
		s.s.Push(v)
	}

	for i := len(values) - 1; i >= 0; i-- {
		v := s.s.Peek()
		s.Equal(values[i], v.(int))

		s.s.Pop()
	}
}

func (s *stackTestSuite) TestLen() {
	s.Equal(0, s.s.Len())

	values := []int{1, 2, 3, 4, 5}
	for i, v := range values {
		s.s.Push(v)
		s.Equal(i+1, s.s.Len())
	}
}

func (s *stackTestSuite) TestEmpty() {
	s.Equal(true, s.s.Empty())

	values := []int{1, 2, 3, 4, 5}
	for _, v := range values {
		s.s.Push(v)
		s.Equal(false, s.s.Empty())
	}
}

func TestStackTestSuite(t *testing.T) {
	s := &stackTestSuite{}
	suite.Run(t, s)
}
