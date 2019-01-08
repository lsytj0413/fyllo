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
	"fmt"
)

// CircularQueue is interface for CircularQueue Data Structure
type CircularQueue interface {
	Empty() bool
	IsFull() bool
	Len() int
	Cap() int
	Poll() interface{}
	Peek() interface{}
	Push(interface{}) bool
}

type circularQueue struct {
	datas []interface{}

	start int
	end   int
	size  int
}

func (c *circularQueue) Push(v interface{}) bool {
	if c.IsFull() {
		return false
	}

	c.size++
	c.datas[c.end] = v
	c.end++
	if c.end >= len(c.datas) {
		c.end = 0
	}
	return true
}

func (c *circularQueue) Peek() interface{} {
	if c.Empty() {
		return nil
	}

	return c.datas[c.start]
}

func (c *circularQueue) Poll() interface{} {
	if c.Empty() {
		return nil
	}

	c.size--
	v := c.datas[c.start]

	c.start++
	if c.start >= len(c.datas) {
		c.start = 0
	}
	return v
}

func (c *circularQueue) IsFull() bool {
	return c.size == len(c.datas)
}

func (c *circularQueue) Empty() bool {
	return c.size == 0
}

func (c *circularQueue) Cap() int {
	return len(c.datas)
}

func (c *circularQueue) Len() int {
	return c.size
}

const (
	minCircularQueueCap int = 1
	maxCircularQueueCap int = 10000
)

// NewCircularQueue will construct CircularQueue object
func NewCircularQueue(cap int) (CircularQueue, error) {
	if cap < minCircularQueueCap || cap > maxCircularQueueCap {
		return nil, fmt.Errorf("Cap should within [%d, %d]", minCircularQueueCap, maxCircularQueueCap)
	}

	return &circularQueue{
		datas: make([]interface{}, cap, cap),
		start: 0,
		end:   0,
		size:  0,
	}, nil
}
