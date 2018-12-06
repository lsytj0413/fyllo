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

package internal

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type generatorTestSuite struct {
	suite.Suite
}

func (s *generatorTestSuite) TestGetByTag() {
	m := uint64(0)
	g, err := NewGenerator(m)
	s.NoError(err)
	sf := g.(*sfGenerator)

	for i := uint64(0); i < MaxTag; i++ {
		g0, ok0 := sf.mappers[i]
		s.Nil(g0)
		s.False(ok0)
	}

	for i := uint64(0); i < MaxTag; i++ {
		g0 := sf.getGeneratorByTag(i)
		s.NotNil(g0)

		g1, ok0 := sf.mappers[i]
		s.Equal(g1, g0)
		s.True(ok0)

		g2 := g0.(*tagGenerator)
		s.Equal(i, g2.tag)
		s.Equal(m, g2.machine)
	}
}

func (s *generatorTestSuite) TestNext() {
	m := uint64(0)
	g, err := NewGenerator(m)
	s.NoError(err)

	now := nano() / Nano2MicroRatio
	tag := uint64(33)
	r, err := g.Next(tag)
	s.NoError(err)
	s.NotNil(r)
	s.Equal(tag, r.Tag)
	s.Equal(m, r.Machine)
	s.True(r.Timestamp >= now)
	s.Equal(uint64(0), r.Sequence)
}

func TestGeneratorTestSuite(t *testing.T) {
	p := &generatorTestSuite{}
	suite.Run(t, p)
}
