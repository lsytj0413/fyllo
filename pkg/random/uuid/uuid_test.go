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

package uuid

import (
	"testing"

	"github.com/lsytj0413/ena/algo/hash"
	"github.com/stretchr/testify/suite"

	"github.com/lsytj0413/fyllo/pkg/random"
)

type uuidProviderTestSuite struct {
	suite.Suite
}

func (s *uuidProviderTestSuite) TestNextOk() {
	p, err := NewProvider(nil)
	s.NoError(err)

	r, err := p.Next(nil)
	s.NoError(err)

	s.Equal(r.Name, ProviderName)
	s.Equal(r.Next, hash.Uint64([]byte(r.Labels[random.LabelIdentify])))
}

func (s *uuidProviderTestSuite) TestNoDuplicate() {
	n := 1000
	store := make(map[uint64]bool, n)
	p, _ := NewProvider(nil)

	for i := 0; i < n; i++ {
		r, _ := p.Next(nil)
		if store[r.Next] {
			s.FailNowf("duplicate next", "value %v", r.Next)
		}
		store[r.Next] = true
	}
}

func TestUuidProviderTestSuite(t *testing.T) {
	s := &uuidProviderTestSuite{}
	suite.Run(t, s)
}
