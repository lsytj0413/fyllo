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
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

type consistentTestSuite struct {
	suite.Suite

	c *consistent
}

type intConsistentHasher struct {
}

func (h *intConsistentHasher) Uint32(data []byte) uint32 {
	v, _ := strconv.Atoi(strings.Replace(string(data), "#", "", -1))
	return uint32(v)
}

func (h *intConsistentHasher) Uint64(data []byte) uint64 {
	v, _ := strconv.Atoi(strings.Replace(string(data), "#", "", -1))
	return uint64(v)
}

func (s *consistentTestSuite) SetupSuite() {
	h := NewHasher(&Config{Replicas: 2})
	s.c = h.(*consistent)
	s.c.hasher = &intConsistentHasher{}

	keys := []string{"1", "3", "2"}
	for i, key := range keys {
		err := s.c.AddNode(NewNode(key, uint(i+1), nil))
		s.NoError(err)
	}
}

func (s *consistentTestSuite) TestHash() {
	values := []string{"11", "12", "14", "22", "27", "34", "35", "29", "19", "27", "100"}
	keys := []string{"1", "2", "2", "2", "3", "1", "1", "3", "2", "3", "1"}

	for i := 0; i < len(values); i++ {
		n, err := s.c.Hash(values[i])
		s.NoError(err)

		s.Equal(keys[i], n.Key())
	}
}

func TestConsistentTestSuite(t *testing.T) {
	s := &consistentTestSuite{}
	suite.Run(t, s)
}
