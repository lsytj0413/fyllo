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

type dequeTestSuite struct {
	suite.Suite

	q Deque
}

func (s *dequeTestSuite) SetupTest() {
	s.q = NewDeque()
}

func (s *dequeTestSuite) TestPushFront() {
	values := []int{1, 2, 3, 4, 5}
	for _, value := range values {
		s.q.PushFront(value)
		s.Equal(value, s.q.PeekFront().(int))
	}
}

func (s *dequeTestSuite) TestPushBack() {
	values := []int{1, 2, 3, 4, 5}
	for _, value := range values {
		s.q.PushBack(value)
		s.Equal(value, s.q.PeekBack().(int))
	}
}

func (s *dequeTestSuite) TestPeekBackNil() {
	v := s.q.PeekBack()
	s.Nil(v)

	s.q.PushBack(0)
	v = s.q.PeekBack()
	s.Equal(0, v.(int))
}

func (s *dequeTestSuite) TestPeekBackOk() {
	values := []int{1, 2, 3, 4, 5}
	targets := []int{1, 1, 3, 3, 5}

	for i := 0; i < len(values); i++ {
		if i%2 == 0 {
			s.q.PushBack(values[i])
		} else {
			s.q.PushFront(values[i])
		}

		v := s.q.PeekBack()
		s.Equal(targets[i], v.(int))
	}
}

func (s *dequeTestSuite) TestPeekFrontNil() {
	v := s.q.PeekFront()
	s.Nil(v)

	s.q.PushBack(0)
	v = s.q.PeekFront()
	s.Equal(0, v.(int))
}

func (s *dequeTestSuite) TestPeekFrontOk() {
	values := []int{1, 2, 3, 4, 5}
	targets := []int{1, 2, 2, 4, 4}

	for i := 0; i < len(values); i++ {
		if i%2 == 0 {
			s.q.PushBack(values[i])
		} else {
			s.q.PushFront(values[i])
		}

		v := s.q.PeekFront()
		s.Equal(targets[i], v.(int))
	}
}

func (s *dequeTestSuite) TestPollBackNil() {
	v := s.q.PollBack()
	s.Equal(nil, v)

	s.q.PushBack(1)
	s.q.PushFront(2)

	v = s.q.PollBack()
	s.Equal(1, v.(int))
	v = s.q.PollBack()
	s.Equal(2, v.(int))
}

func (s *dequeTestSuite) TestPollBackOk() {
	values := []int{1, 2, 3, 4, 5}
	for i, value := range values {
		if i%2 == 0 {
			s.q.PushBack(value)
		} else {
			s.q.PushFront(value)
		}
	}

	targets := []int{5, 3, 1, 2, 4}
	for _, target := range targets {
		v := s.q.PollBack()
		s.Equal(target, v.(int))
	}

	s.Nil(s.q.PollBack())
}

func (s *dequeTestSuite) TestPollFrontNil() {
	v := s.q.PollFront()
	s.Equal(nil, v)

	s.q.PushBack(1)
	s.q.PushFront(2)

	v = s.q.PollFront()
	s.Equal(2, v.(int))
	v = s.q.PollFront()
	s.Equal(1, v.(int))
}

func (s *dequeTestSuite) TestPollFrontOk() {
	values := []int{1, 2, 3, 4, 5}
	for i, value := range values {
		if i%2 == 0 {
			s.q.PushBack(value)
		} else {
			s.q.PushFront(value)
		}
	}

	targets := []int{4, 2, 1, 3, 5}
	for _, target := range targets {
		v := s.q.PollFront()
		s.Equal(target, v.(int))
	}

	s.Nil(s.q.PollFront())
}

func (s *dequeTestSuite) TestEmpty() {
	s.Equal(true, s.q.Empty())

	values := []int{1, 2, 3, 4, 5}
	for i, value := range values {
		if i%2 == 0 {
			s.q.PushBack(value)
		} else {
			s.q.PushFront(value)
		}

		s.Equal(false, s.q.Empty())
	}
}

func (s *dequeTestSuite) TestLen() {
	s.Equal(0, s.q.Len())

	values := []int{1, 2, 3, 4, 5}
	for i, value := range values {
		if i%2 == 0 {
			s.q.PushBack(value)
		} else {
			s.q.PushFront(value)
		}

		s.Equal(i+1, s.q.Len())
	}
}

func TestDequeTestSuite(t *testing.T) {
	s := &dequeTestSuite{}
	suite.Run(t, s)
}
