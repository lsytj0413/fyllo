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
	"sync"
	"testing"

	"github.com/stretchr/testify/suite"
)

type consistentSafeTestSuite struct {
	suite.Suite

	c *safeConsistent
}

func (s *consistentSafeTestSuite) SetupSuite() {
	h := NewSafeHasher(&Config{Replicas: 2})
	s.c = h.(*safeConsistent)

	c := s.c.Hasher.(*consistent)
	c.hasher = &intConsistentHasher{}

	keys := []string{"1", "3", "2"}
	for i, key := range keys {
		err := s.c.AddNode(NewNode(key, uint(i+1), nil))
		s.NoError(err)
	}
}

func (s *consistentSafeTestSuite) TestHash() {
	values := []string{"11", "12", "14", "22", "27", "34", "35", "29", "19", "27", "100"}
	keys := []string{"1", "2", "2", "2", "3", "1", "1", "3", "2", "3", "1"}

	count := 10
	var wg sync.WaitGroup

	for i := 0; i < count; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for i := 0; i < len(values); i++ {
				n, err := s.c.Hash(values[i])
				s.NoError(err)

				s.Equal(keys[i], n.Key())
			}
		}()
	}

	wg.Wait()
}
func TestConsistentSafeTestSuite(t *testing.T) {
	s := &consistentSafeTestSuite{}
	suite.Run(t, s)
}
