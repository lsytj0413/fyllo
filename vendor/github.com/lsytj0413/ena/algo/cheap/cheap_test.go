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

package cheap

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

type cheapTestSuite struct {
	suite.Suite

	h Heap
}

type item struct {
	id       uint64
	priority uint64
}

func less(v1 Interface, v2 Interface) bool {
	i1 := v1.(*item)
	i2 := v2.(*item)

	return i1.priority < i2.priority
}

func (i *item) Id() uint64 {
	return i.id
}

func newItem(id uint64, priority uint64) *item {
	return &item{
		id:       id,
		priority: priority,
	}
}

func (s *cheapTestSuite) SetupTest() {
	s.h = NewHeap(less)
}

func (s *cheapTestSuite) TestTopOk() {
	s.h.PushX(newItem(1, 2))
	s.h.PushX(newItem(0, 0))

	x := (s.h.Top()).(Interface)
	s.Equal(uint64(0), x.Id())
}

func (s *cheapTestSuite) TestTopNil() {
	x := s.h.Top()
	s.Nil(x)
}

func (s *cheapTestSuite) TestPopXOk() {
	s.h.PushX(newItem(1, 2))
	s.h.PushX(newItem(0, 0))

	x := (s.h.PopX()).(Interface)
	s.Equal(uint64(0), x.Id())

	x = (s.h.PopX()).(Interface)
	s.Equal(uint64(1), x.Id())
}

func (s *cheapTestSuite) TestPopXNil() {
	x := s.h.PopX()
	s.Nil(x)
}

func (s *cheapTestSuite) TestPushXOk() {
	err := s.h.PushX(newItem(1, 2))
	s.NoError(err)
}

func (s *cheapTestSuite) TestPushXExists() {
	i := newItem(1, 2)
	err := s.h.PushX(i)
	s.NoError(err)

	err = s.h.PushX(i)

	s.Error(err, fmt.Errorf("Heap.PushX: id=%d exists", i.Id()).Error())
}

func (s *cheapTestSuite) TestUpdateOk() {
	s.h.PushX(newItem(1, 2))
	s.h.PushX(newItem(0, 0))

	x := (s.h.Top()).(*item)
	x.priority = 3

	err := s.h.Update(x)
	s.Nil(err)

	y := (s.h.Top()).(Interface)
	s.Equal(uint64(1), y.Id())
}

func (s *cheapTestSuite) TestUpdateFail() {
	i := newItem(0, 3)
	err := s.h.Update(i)
	s.EqualError(err, fmt.Errorf("Heap.Update: id=%d not exists", i.Id()).Error())

	s.h.PushX(newItem(1, 2))
	err = s.h.Update(i)
	s.EqualError(err, fmt.Errorf("Heap.Update: id=%d not exists", i.Id()).Error())
}

func (s *cheapTestSuite) TestRemoveOK() {
	s.h.PushX(newItem(1, 2))
	s.h.PushX(newItem(0, 0))

	x := (s.h.Top()).(*item)

	err := s.h.Remove(x)
	s.NoError(err)

	y := (s.h.Top()).(Interface)
	s.Equal(uint64(1), y.Id())
}

func (s *cheapTestSuite) TestRemoveNotExists() {
	i := newItem(0, 3)
	err := s.h.Remove(i)
	s.EqualError(err, fmt.Errorf("Heap.Remove: id=%d not exists", i.Id()).Error())
}

func TestCheapTestSuite(t *testing.T) {
	s := &cheapTestSuite{}
	suite.Run(t, s)
}
