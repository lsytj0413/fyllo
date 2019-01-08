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

// Queue is interface for Queue Data Structure
type Queue interface {
	Empty() bool
	Len() int
	Poll() interface{}
	Peek() interface{}
	Push(interface{})
}

type queue struct {
	list *list.List
}

// NewQueue construct Queue
func NewQueue() Queue {
	return &queue{
		list: list.New(),
	}
}

func (q *queue) Empty() bool {
	return q.list.Len() == 0
}

func (q *queue) Len() int {
	return q.list.Len()
}

func (q *queue) Poll() interface{} {
	e := q.list.Front()
	if e != nil {
		q.list.Remove(e)
		return e.Value
	}

	return nil
}

func (q *queue) Peek() interface{} {
	e := q.list.Front()
	if e != nil {
		return e.Value
	}

	return nil
}

func (q *queue) Push(v interface{}) {
	q.list.PushBack(v)
}
