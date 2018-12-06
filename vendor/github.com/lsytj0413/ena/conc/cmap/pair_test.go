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

type pairTestSuite struct {
	suite.Suite
}

func (s *pairTestSuite) TestNewPairOk() {
	key := "key"
	p, err := newPair(key, key)

	s.NoError(err)
	s.NotNil(p)

	s.Equal(key, p.Key())
	v := p.Value().(string)
	s.Equal(key, v)
}

func (s *pairTestSuite) TestNewPairValueNil() {
	p, err := newPair("key", nil)
	s.Nil(p)

	s.Error(err)
}

func (s *pairTestSuite) TestValueNil() {
	p := &pair{}

	v := p.Value()
	s.Nil(v)
}

func (s *pairTestSuite) TestSetValueOk() {
	p := &pair{}

	v0 := "123"
	err := p.SetValue(v0)

	s.NoError(err)

	v := p.Value().(string)
	s.Equal(v0, v)
}

func (s *pairTestSuite) TestSetValueNil() {
	p := &pair{}

	err := p.SetValue(nil)
	s.Error(err)
}

func (s *pairTestSuite) TestNextOk() {
	p := &pair{}
	n := &pair{}

	p.SetNext(n)
	v := p.Next()

	s.NotNil(v)
	s.Equal(n, v)
}

func (s *pairTestSuite) TestNextNil() {
	p := &pair{}

	v := p.Next()
	s.Nil(v)
}

func (s *pairTestSuite) TestSetNextOk() {
	p := &pair{}
	n := &pair{}

	p.SetNext(n)
	v := p.Next()

	s.NotNil(v)
	s.Equal(n, v)
}

func (s *pairTestSuite) TestSetNextNil() {
	p := &pair{}

	err := p.SetNext(nil)
	s.NoError(err)
}

type fackPair struct {
}

func (*fackPair) Key() string {
	return ""
}

func (*fackPair) Hash() uint64 {
	return 0
}

func (*fackPair) Value() interface{} {
	return nil
}

func (*fackPair) SetValue(v interface{}) error {
	return nil
}

func (*fackPair) Clone() Pair {
	return nil
}

func (*fackPair) String() string {
	return ""
}

func (*fackPair) Next() Pair {
	return nil
}

func (*fackPair) SetNext(n Pair) error {
	return nil
}

func (s *pairTestSuite) TestSetNextInvalidPairType() {
	p := &pair{}

	err := p.SetNext(&fackPair{})
	s.Error(err)
}

func (s *pairTestSuite) TestClone() {
	p, _ := newPair("key", "key")
	c := p.Clone()

	s.Equal(p.Key(), c.Key())
	s.Equal(p.Value().(string), c.Value().(string))
}

func TestPairTestSuite(t *testing.T) {
	s := &pairTestSuite{}
	suite.Run(t, s)
}
