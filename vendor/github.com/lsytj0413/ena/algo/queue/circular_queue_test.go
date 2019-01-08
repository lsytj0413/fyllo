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

type circularQueueTestSuite struct {
	suite.Suite

	q *circularQueue
}

var (
	testCap = 3
)

func (s *circularQueueTestSuite) SetupTest() {
	v, err := NewCircularQueue(testCap)
	s.NoError(err)

	s.q = v.(*circularQueue)
}

func (s *circularQueueTestSuite) TestCap() {
	s.Equal(testCap, s.q.Cap())
}

func (s *circularQueueTestSuite) TestEmpty() {
	s.True(s.q.Empty())

	values := []int{1, 2, 3}
	for _, v := range values {
		s.q.Push(v)
		s.False(s.q.Empty())
	}
}

func (s *circularQueueTestSuite) TestLen() {
	s.Equal(0, s.q.Len())

	values := []int{1, 2, 3}
	for i, v := range values {
		s.q.Push(v)
		s.Equal(i+1, s.q.Len())
	}
}

func (s *circularQueueTestSuite) TestIsFull() {
	s.False(s.q.IsFull())

	values := []int{1, 2}
	for _, v := range values {
		s.q.Push(v)
		s.False(s.q.IsFull())
	}

	s.q.Push(3)
	s.True(s.q.IsFull())
}

func (s *circularQueueTestSuite) TestPush() {
	values := []int{1, 2, 3}
	for _, v := range values {
		s.True(s.q.Push(v))
	}

	values = []int{4, 5, 6, 7}
	for _, v := range values {
		s.False(s.q.Push(v))
	}

	s.q.Poll()
	s.q.Poll()

	values = []int{4, 5}
	for _, v := range values {
		s.True(s.q.Push(v))
	}
}

func (s *circularQueueTestSuite) TestPoll() {
	s.Nil(s.q.Poll())

	s.q.Push(1)
	s.q.Push(2)
	s.q.Push(3)

	s.Equal(1, s.q.Poll().(int))
	s.Equal(2, s.q.Poll().(int))

	s.q.Push(4)
	s.Equal(3, s.q.Poll().(int))
	s.Equal(4, s.q.Poll().(int))

	s.Nil(s.q.Poll())
}

func (s *circularQueueTestSuite) TestPeek() {
	s.Nil(s.q.Peek())

	s.q.Push(1)
	s.q.Push(2)
	s.q.Push(3)
	s.Equal(1, s.q.Peek().(int))
	s.q.Poll()

	s.Equal(2, s.q.Peek().(int))
	s.q.Poll()

	s.q.Push(4)
	s.Equal(3, s.q.Peek().(int))
	s.q.Poll()

	s.Equal(4, s.q.Peek().(int))
	s.q.Poll()

	s.Nil(s.q.Peek())
}

func TestCircularQueueTestSuite(t *testing.T) {
	s := &circularQueueTestSuite{}
	suite.Run(t, s)
}
