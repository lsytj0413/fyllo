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

import "fmt"

// CircularDeque is interface for CircularDeque Data Structure
type CircularDeque interface {
	Empty() bool
	IsFull() bool
	Len() int
	Cap() int
	PollFront() interface{}
	PollBack() interface{}
	PeekFront() interface{}
	PeekBack() interface{}
	PushFront(interface{}) bool
	PushBack(interface{}) bool
}

type circularDeque struct {
	datas []interface{}

	start int
	end   int
	size  int
}

func (c *circularDeque) PushFront(v interface{}) bool {
	if c.IsFull() {
		return false
	}

	c.size++
	c.start--
	if c.start < 0 {
		c.start = len(c.datas) - 1
	}
	c.datas[c.start] = v
	return true
}

func (c *circularDeque) PushBack(v interface{}) bool {
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

func (c *circularDeque) PeekBack() interface{} {
	if c.Empty() {
		return nil
	}

	if c.end == 0 {
		return c.datas[len(c.datas)-1]
	}
	return c.datas[c.end-1]
}

func (c *circularDeque) PeekFront() interface{} {
	if c.Empty() {
		return nil
	}

	return c.datas[c.start]
}

func (c *circularDeque) PollBack() interface{} {
	if c.Empty() {
		return nil
	}

	c.size--
	if c.end == 0 {
		c.end = len(c.datas) - 1
	} else {
		c.end--
	}
	v := c.datas[c.end]
	c.datas[c.end] = nil

	return v
}

func (c *circularDeque) PollFront() interface{} {
	if c.Empty() {
		return nil
	}

	c.size--
	v := c.datas[c.start]
	c.datas[c.start] = nil

	c.start++
	if c.start >= len(c.datas) {
		c.start = 0
	}
	return v
}

func (c *circularDeque) IsFull() bool {
	return c.size == len(c.datas)
}

func (c *circularDeque) Empty() bool {
	return c.size == 0
}

func (c *circularDeque) Cap() int {
	return len(c.datas)
}

func (c *circularDeque) Len() int {
	return c.size
}

const (
	minCircularDequeCap int = 1
	maxCircularDequeCap int = 10000
)

// NewCircularDeque will construct CircularDeque object
func NewCircularDeque(cap int) (CircularDeque, error) {
	if cap < minCircularDequeCap || cap > maxCircularDequeCap {
		return nil, fmt.Errorf("Cap should within [%d, %d]", minCircularDequeCap, maxCircularDequeCap)
	}

	return &circularDeque{
		datas: make([]interface{}, cap, cap),
		start: 0,
		end:   0,
		size:  0,
	}, nil
}
