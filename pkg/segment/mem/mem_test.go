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

package mem

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/lsytj0413/fyllo/pkg/segment"
	"github.com/lsytj0413/fyllo/pkg/segment/internal"
)

type memTestSuite struct {
	suite.Suite
}

func (s *memTestSuite) TestProviderNextOk() {
	tagName := "test"
	m := &memStorage{
		tags: map[string]*memRow{
			tagName: {
				Tag:         tagName,
				Max:         0,
				Step:        1,
				Description: "test row",
			},
		},
	}
	p, err := internal.NewProvider(ProviderName, m)
	s.NoError(err)
	s.NotNil(p)

	var r *segment.Result
	for i := 0; i < 100; i++ {
		r, err = p.Next(&segment.Arguments{
			Tag: tagName,
		})
		s.NoError(err)
		s.NotNil(r)

		s.Equal(uint64(i), r.Next)
		s.Equal(ProviderName, r.Name)
		s.Equal(1, len(r.Labels))
		s.Equal(tagName, r.Labels[segment.LabelTag])
	}
}

func TestMemTestSuite(t *testing.T) {
	s := &memTestSuite{}
	suite.Run(t, s)
}
