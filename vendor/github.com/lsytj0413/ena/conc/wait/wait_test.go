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

package wait

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

type waitTestSuite struct {
	suite.Suite
	w Wait
}

func (s *waitTestSuite) SetupTest() {
	s.w = New()
}

func (s *waitTestSuite) TearDownTest() {
	s.w = nil
}

func (s *waitTestSuite) TestNew() {
	w := New()

	_, ok := w.(*defWait)
	s.True(ok)
}

func (s *waitTestSuite) TestIsRegisteredOk() {
	for i := uint64(0); i < uint64(100); i++ {
		_, err := s.w.Register(i)
		s.NoError(err)
	}

	for i := uint64(0); i < uint64(100); i++ {
		ok := s.w.IsRegistered(i)
		s.True(ok)
	}
}

func (s *waitTestSuite) TestIsRegisterdFail() {
	for i := uint64(0); i < uint64(100); i++ {
		ok := s.w.IsRegistered(i)
		s.False(ok)
	}
}

func (s *waitTestSuite) TestRegisterOk() {
	indexs := []uint64{0, 99, 100, 234234, 34535435345345, ^uint64(0)}
	for _, i := range indexs {
		_, err := s.w.Register(i)
		s.NoError(err)
	}
}

func (s *waitTestSuite) TestRegisterDupId() {
	indexs := []uint64{0, 99, 100, 234234, 34535435345345, ^uint64(0)}
	for _, i := range indexs {
		_, err := s.w.Register(i)
		s.NoError(err)
	}

	for _, i := range indexs {
		_, err := s.w.Register(i)
		s.Error(err)
	}
}

func (s *waitTestSuite) TestTriggerOk() {
	id := uint64(98989)
	ch, err := s.w.Register(id)
	s.NoError(err)

	err = s.w.Trigger(id, id)
	s.NoError(err)

	// check value
	v := <-ch
	v0, ok := v.(uint64)
	s.True(ok)
	s.Equal(id, v0)

	// check register
	ok = s.w.IsRegistered(id)
	s.False(ok)
}

func (s *waitTestSuite) TestTriggerNotRegister() {
	id := uint64(98989)
	err := s.w.Trigger(id, id)
	s.Error(err)

	s.Equal(fmt.Sprintf("Wait.Trigger Failed: id=%d not registered", id), err.Error())
}

func (s *waitTestSuite) TestTriggerTimeout() {
	id := uint64(34324)
	_, _ = s.w.Register(id)

	// chn, _ := ch.(chan interface{})
	// chn <- nil
	wn := s.w.(*defWait)
	ch, _ := wn.m[id]
	ch <- nil

	err := s.w.Trigger(id, id)
	s.Error(err)
	s.Equal(fmt.Sprintf("Wait.Trigger Failed: id=%d timeout", id), err.Error())
}

func TestWaitTestSuite(t *testing.T) {
	s := &waitTestSuite{}
	suite.Run(t, s)
}
