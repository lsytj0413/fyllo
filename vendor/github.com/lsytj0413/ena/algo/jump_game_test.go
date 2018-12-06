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

package algo

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type jumpgameTestSuite struct {
	suite.Suite
}

type jump1 struct {
	a1 []uint
	r  bool
}

func (s *jumpgameTestSuite) TestJump1() {
	vs := []jump1{
		{
			a1: []uint{1, 4, 1, 1, 4},
			r:  true,
		},
		{
			a1: []uint{3, 2, 1, 0, 4},
			r:  false,
		},
		{
			a1: []uint{},
			r:  true,
		},
		{
			a1: []uint{0},
			r:  true,
		},
		{
			a1: []uint{0, 2},
			r:  false,
		},
		{
			a1: []uint{2, 3, 0, 0, 0},
			r:  true,
		},
		{
			a1: []uint{2, 3, 0, 0, 0, 1},
			r:  false,
		},
	}

	for _, v := range vs {
		s.Equal(v.r, IsJumpEnd(v.a1))
	}
}

type jump2 struct {
	a1 []uint
	r  int
}

func (s *jumpgameTestSuite) TestJump2() {
	vs := []jump2{
		{
			a1: []uint{1, 4, 1, 1, 4},
			r:  2,
		},
		{
			a1: []uint{2, 3, 1, 1, 4},
			r:  2,
		},
		{
			a1: []uint{2, 5, 1, 1, 4},
			r:  2,
		},
		{
			a1: []uint{},
			r:  -1,
		},
		{
			a1: []uint{1},
			r:  0,
		},
	}

	for _, v := range vs {
		s.Equal(v.r, JumpShortedSteps(v.a1))
	}
}
func TestJumpgameTestSuite(t *testing.T) {
	s := &jumpgameTestSuite{}
	suite.Run(t, s)
}
