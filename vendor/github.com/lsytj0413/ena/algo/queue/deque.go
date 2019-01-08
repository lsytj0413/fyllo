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
	"container/list"
)

// Deque is interface for Deque Data Structure
type Deque interface {
	Empty() bool
	Len() int
	PollFront() interface{}
	PollBack() interface{}
	PeekFront() interface{}
	PeekBack() interface{}
	PushFront(interface{})
	PushBack(interface{})
}

type deque struct {
	l *list.List
}

func (q *deque) PushFront(v interface{}) {
	q.l.PushFront(v)
}

func (q *deque) PushBack(v interface{}) {
	q.l.PushBack(v)
}

func (q *deque) PeekBack() interface{} {
	e := q.l.Back()
	if e == nil {
		return nil
	}

	return e.Value
}

func (q *deque) PeekFront() interface{} {
	e := q.l.Front()
	if e == nil {
		return nil
	}

	return e.Value
}

func (q *deque) PollBack() interface{} {
	e := q.l.Back()
	if e == nil {
		return nil
	}

	q.l.Remove(e)
	return e.Value
}

func (q *deque) PollFront() interface{} {
	e := q.l.Front()
	if e == nil {
		return nil
	}

	q.l.Remove(e)
	return e.Value
}

func (q *deque) Empty() bool {
	return q.l.Len() == 0
}

func (q *deque) Len() int {
	return q.l.Len()
}

// NewDeque construct Deque
func NewDeque() Deque {
	return &deque{
		l: list.New(),
	}
}
