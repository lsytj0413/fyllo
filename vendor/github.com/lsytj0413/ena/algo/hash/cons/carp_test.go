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

package cons

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/suite"
)

type carpTestSuite struct {
	suite.Suite

	h *carp
}

type intHasher struct {
}

func (h *intHasher) Uint32(data []byte) uint32 {
	v, _ := strconv.Atoi(string(data))
	return uint32(v)
}

func (h *intHasher) Uint64(data []byte) uint64 {
	v, _ := strconv.Atoi(string(data))
	return uint64(v)
}

func (s *carpTestSuite) SetupSuite() {
	h, _ := NewCarp([]string{"1", "2", "3"})
	s.h = h.(*carp)
	s.h.hasher = &intHasher{}
}

func (s *carpTestSuite) TestHash() {
	keys := []string{"3", "1", "2", "-"}
	expected := []string{"1", "1", "1", "3"}
	for i, key := range keys {
		v, _ := s.h.Hash(key)
		s.Equal(expected[i], v)
	}
}

func (s *carpTestSuite) TestNewCarpError() {
	endpoints := []string{"1", "2", "3"}
	hasError := []bool{true, false, false}
	for i := 0; i < len(endpoints); i++ {
		_, err := NewCarp(endpoints[0:i])
		if hasError[i] {
			s.Error(err)
		} else {
			s.NoError(err)
		}
	}
}

func TestCarpTestSuite(t *testing.T) {
	s := &carpTestSuite{}
	suite.Run(t, s)
}
