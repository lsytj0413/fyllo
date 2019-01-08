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

package hash

import (
	"hash/fnv"
	"sync"
	"testing"

	"github.com/stretchr/testify/suite"
)

type hashTestSuite struct {
	suite.Suite
}

func hashUint32(data []byte) uint32 {
	v32 := fnv.New32()
	v32.Write(data)
	return v32.Sum32()
}

func hashUint64(data []byte) uint64 {
	v64 := fnv.New64()
	v64.Write(data)
	return v64.Sum64()
}

func (s *hashTestSuite) TestHasherUint32() {
	datas := []string{"abc", "def", "xxxxxx", "100", "x.43.}"}
	h := NewHash()

	for _, data := range datas {
		s.Equal(hashUint32([]byte(data)), h.Uint32([]byte(data)))
	}
}

func (s *hashTestSuite) TestHasherUint64() {
	datas := []string{"abc", "def", "xxxxxx", "100", "x.43.}"}
	h := NewHash()

	for _, data := range datas {
		s.Equal(hashUint64([]byte(data)), h.Uint64([]byte(data)))
	}
}

func (s *hashTestSuite) TestSafeHasherUint32() {
	datas := []string{"abc", "def", "xxxxxx", "100", "x.43.}"}
	h := NewSafeHash()

	count := 10
	var wg sync.WaitGroup

	for i := 0; i < count; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for _, data := range datas {
				expect := hashUint32([]byte(data))
				actual := h.Uint32([]byte(data))
				s.Equal(expect, actual)
			}
		}()
	}

	wg.Wait()
}

func (s *hashTestSuite) TestSafeHasherUint64() {
	datas := []string{"abc", "def", "xxxxxx", "100", "x.43.}"}
	h := NewSafeHash()

	count := 10
	var wg sync.WaitGroup

	for i := 0; i < count; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for _, data := range datas {
				expect := hashUint64([]byte(data))
				actual := h.Uint64([]byte(data))
				s.Equal(expect, actual)
			}
		}()
	}

	wg.Wait()
}

func TestHashTestSuite(t *testing.T) {
	s := &hashTestSuite{}
	suite.Run(t, s)
}
