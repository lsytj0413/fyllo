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
package queue

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type queueTestSuite struct {
	suite.Suite

	s Queue
}

func (s *queueTestSuite) SetupTest() {
	s.s = NewQueue()
}

func (s *queueTestSuite) TestPushOk() {
	r := s.s.Peek()
	s.Nil(r)

	values := []int{1, 2, 3, 4, 5}
	for i, v := range values {
		s.s.Push(v)
		s.Equal(i+1, s.s.Len())
	}
}

func (s *queueTestSuite) TestPollNil() {
	r := s.s.Poll()
	s.Nil(r)
}

func (s *queueTestSuite) TestPollOk() {
	values := []int{1, 2, 3, 4, 5}
	for _, v := range values {
		s.s.Push(v)
	}

	for i := 0; i < len(values); i++ {
		v := s.s.Poll()
		s.Equal(values[i], v.(int))
	}

	s.Equal(0, s.s.Len())
}

func (s *queueTestSuite) TestPeekNil() {
	r := s.s.Peek()
	s.Nil(r)
}

func (s *queueTestSuite) TestPeekOk() {
	values := []int{1, 2, 3, 4, 5}
	for _, v := range values {
		s.s.Push(v)
	}

	for i := 0; i < len(values); i++ {
		v := s.s.Peek()
		s.Equal(values[i], v.(int))

		s.s.Poll()
	}
}

func (s *queueTestSuite) TestLen() {
	s.Equal(0, s.s.Len())

	values := []int{1, 2, 3, 4, 5}
	for i, v := range values {
		s.s.Push(v)
		s.Equal(i+1, s.s.Len())
	}
}

func (s *queueTestSuite) TestEmpty() {
	s.Equal(true, s.s.Empty())

	values := []int{1, 2, 3, 4, 5}
	for _, v := range values {
		s.s.Push(v)
		s.Equal(false, s.s.Empty())
	}
}

func TestQueueTestSuite(t *testing.T) {
	s := &queueTestSuite{}
	suite.Run(t, s)
}
