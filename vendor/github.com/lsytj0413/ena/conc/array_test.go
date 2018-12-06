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

package conc

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/suite"
)

type concArrayTestSuite struct {
	suite.Suite
}

func (s *concArrayTestSuite) TestNewArrayOk() {
	var i uint32
	for ; i <= maxArrayLength; i++ {
		a, err := NewConcurrentArray(i)
		s.NotNil(a)
		s.NoError(err)
	}
}

func (s *concArrayTestSuite) TestNewArrayError() {
	var i uint32 = maxArrayLength + 1
	for ; i < maxArrayLength+10; i++ {
		a, err := NewConcurrentArray(i)
		s.Nil(a)
		s.Error(err, fmt.Sprintf("NewConcurrentArray: length must be less than [%v]", maxArrayLength+1))
	}
}

func (s *concArrayTestSuite) TestSetOk() {
	var len uint32 = 50
	a, err := NewConcurrentArray(len)
	s.NoError(err)

	for i := uint32(0); i < len; i++ {
		err = a.Set(i, i)
		s.NoError(err)
	}
}

func (s *concArrayTestSuite) TestSetError() {
	var len uint32 = 50
	a, err := NewConcurrentArray(len)
	s.NoError(err)

	for i := len; i < len+10; i++ {
		err = a.Set(i, i)
		s.Error(err, fmt.Sprintf("ConcurrentArray.checkIndex: Index out of range [0, %d)", a.Len()))
	}
}

func (s *concArrayTestSuite) TestGetOk() {
	var len uint32 = 50
	a, err := NewConcurrentArray(len)
	s.NoError(err)

	for i := uint32(0); i < len; i++ {
		err = a.Set(i, i)
		s.NoError(err)
	}

	for i := uint32(0); i < len; i++ {
		v, err := a.Get(i)
		s.NoError(err)
		s.Equal(i, v.(uint32))
	}
}

func (s *concArrayTestSuite) TestGetError() {
	var len uint32 = 50
	a, err := NewConcurrentArray(len)
	s.NoError(err)

	for i := len; i < len+10; i++ {
		v, err := a.Get(i)
		s.Error(err, fmt.Sprintf("ConcurrentArray.checkIndex: Index out of range [0, %d)", a.Len()))
		s.Nil(v)
	}
}

func (s *concArrayTestSuite) TestLen() {
	var i uint32
	for ; i <= maxArrayLength; i++ {
		a, err := NewConcurrentArray(i)
		s.NotNil(a)
		s.NoError(err)

		s.Equal(i, a.Len())
	}
}

// test set for multi goroutine
func (s *concArrayTestSuite) TestConcSetOk() {
	arr, err := NewConcurrentArray(10)
	s.NoError(err)
	var wg sync.WaitGroup
	wg.Add(int(arr.Len()))
	max := 1000
	for i := 0; i < int(arr.Len()); i++ {
		go func(i int) {
			defer wg.Done()

			for j := 0; j < max; j++ {
				err := arr.Set(uint32(i), j)
				s.NoError(err)
			}
		}(i)
	}

	wg.Wait()
	for i := uint32(0); i < arr.Len(); i++ {
		v, err := arr.Get(i)
		s.NoError(err)
		s.Equal(max-1, v.(int))
	}
}

func TestConcArrayTestSuite(t *testing.T) {
	s := &concArrayTestSuite{}
	suite.Run(t, s)
}
