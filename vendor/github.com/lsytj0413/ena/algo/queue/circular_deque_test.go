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

type circularDequeTestSuite struct {
	suite.Suite

	q *circularDeque
}

func (s *circularDequeTestSuite) SetupTest() {
	v, err := NewCircularDeque(testCap)
	s.NoError(err)

	s.q = v.(*circularDeque)
}

func (s *circularDequeTestSuite) TestCap() {
	s.Equal(testCap, s.q.Cap())
}

func (s *circularDequeTestSuite) TestEmpty() {
	s.True(s.q.Empty())

	values := []int{1, 2, 3}
	for i, v := range values {
		if i%2 == 0 {
			s.q.PushBack(v)
		} else {
			s.q.PushFront(v)
		}
		s.False(s.q.Empty())
	}
}

func (s *circularDequeTestSuite) TestLen() {
	s.Equal(0, s.q.Len())

	values := []int{1, 2, 3}
	for i, v := range values {
		if i%2 == 0 {
			s.q.PushBack(v)
		} else {
			s.q.PushFront(v)
		}
		s.Equal(i+1, s.q.Len())
	}
}

func (s *circularDequeTestSuite) TestIsFull() {
	s.False(s.q.IsFull())

	values := []int{1, 2}
	for i, v := range values {
		if i%2 == 0 {
			s.q.PushBack(v)
		} else {
			s.q.PushFront(v)
		}
		s.False(s.q.IsFull())
	}

	s.q.PushBack(3)
	s.True(s.q.IsFull())
}

func (s *circularDequeTestSuite) TestPush() {
	values := []int{1, 2, 3}
	for _, v := range values {
		s.True(s.q.PushBack(v))
	}

	values = []int{4, 5, 6, 7}
	for _, v := range values {
		s.False(s.q.PushBack(v))
	}
	for _, v := range values {
		s.False(s.q.PushFront(v))
	}

	s.q.PollFront()
	s.q.PollFront()

	values = []int{4}
	for _, v := range values {
		s.True(s.q.PushBack(v))
	}

	values = []int{5}
	for _, v := range values {
		s.True(s.q.PushFront(v))
	}
}

func (s *circularDequeTestSuite) TestPoll() {
	s.Nil(s.q.PollFront())

	s.q.PushBack(1)
	s.q.PushFront(2)
	s.q.PushBack(3)

	s.Equal(2, s.q.PollFront().(int))
	s.Equal(3, s.q.PollBack().(int))

	s.q.PushBack(4)
	s.Equal(4, s.q.PollBack().(int))
	s.Equal(1, s.q.PollFront().(int))

	s.Nil(s.q.PollFront())
	s.Nil(s.q.PollBack())
}

func (s *circularDequeTestSuite) TestPeek() {
	s.Nil(s.q.PeekFront())

	s.q.PushBack(1)
	s.q.PushFront(2)
	s.q.PushBack(3)

	s.Equal(2, s.q.PeekFront().(int))
	s.Equal(3, s.q.PeekBack().(int))
	s.q.PollBack()

	s.q.PushBack(4)
	s.Equal(4, s.q.PeekBack().(int))
	s.Equal(2, s.q.PeekFront().(int))

	s.q.PollBack()
	s.q.PollFront()

	s.Equal(1, s.q.PeekBack().(int))
	s.Equal(1, s.q.PeekFront().(int))
}

func TestCircularDequeTestSuite(t *testing.T) {
	s := &circularDequeTestSuite{}
	suite.Run(t, s)
}
