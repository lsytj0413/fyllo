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

package cmap

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type mapTestSuite struct {
	suite.Suite
}

func (s *mapTestSuite) TestPutOk() {
	m, err := NewMap(uint32(10))
	s.NoError(err)

	b, err := m.Put("key", "value")
	s.NoError(err)
	s.True(b)
}

func (s *mapTestSuite) TestGetOk() {
	m, err := NewMap(uint32(10))
	s.NoError(err)

	b, err := m.Put("key", "value")
	s.NoError(err)
	s.True(b)

	v, b := m.Get("key")
	s.True(b)
	s.Equal("value", v.(string))

	b, err = m.Put("key", "value2")
	s.NoError(err)
	s.False(b)

	v, b = m.Get("key")
	s.True(b)
	s.Equal("value2", v.(string))
}

func TestMapTestSuite(t *testing.T) {
	s := &mapTestSuite{}
	suite.Run(t, s)
}
