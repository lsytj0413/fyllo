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

import "container/list"

// Stack is interface for Stack Data Structure
type Stack interface {
	Push(interface{})
	Pop() interface{}
	Peek() interface{}
	Len() int
	Empty() bool
}

type stack struct {
	list *list.List
}

// New construct Stack
func New() Stack {
	return &stack{
		list: list.New(),
	}
}

func (s *stack) Push(v interface{}) {
	s.list.PushBack(v)
}

func (s *stack) Pop() interface{} {
	e := s.list.Back()
	if e != nil {
		s.list.Remove(e)
		return e.Value
	}

	return nil
}

func (s *stack) Peek() interface{} {
	e := s.list.Back()
	if e != nil {
		return e.Value
	}

	return nil
}

func (s *stack) Len() int {
	return s.list.Len()
}

func (s *stack) Empty() bool {
	return s.list.Len() == 0
}
