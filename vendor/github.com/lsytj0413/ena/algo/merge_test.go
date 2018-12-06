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

type mergeTestSuite struct {
	suite.Suite
}

type merge struct {
	a1 []int
	a2 []int
	r  []int
}

func (s *mergeTestSuite) TestEmpty() {
	v := merge{
		a1: make([]int, 0),
		a2: make([]int, 0),
		r:  make([]int, 0),
	}

	r := Merge(v.a1, v.a2)
	s.Equal(v.r, r)
}

func (s *mergeTestSuite) TestOne() {
	vs := []merge{
		{
			a1: make([]int, 0),
			a2: []int{0},
			r:  []int{0},
		},
		{
			a1: make([]int, 0),
			a2: []int{0, 2, 3},
			r:  []int{0, 2, 3},
		},
		{
			a1: []int{4},
			a2: []int{},
			r:  []int{4},
		},
		{
			a1: []int{0, 4, 6},
			a2: []int{},
			r:  []int{0, 4, 6},
		},
	}

	for _, v := range vs {
		r := Merge(v.a1, v.a2)
		s.Equal(v.r, r)
	}
}

func (s *mergeTestSuite) TestTwo() {
	vs := []merge{
		{
			a1: []int{0, 2, 3},
			a2: []int{0},
			r:  []int{0, 0, 2, 3},
		},
		{
			a1: []int{0, 4, 7},
			a2: []int{0, 2, 3},
			r:  []int{0, 0, 2, 3, 4, 7},
		},
		{
			a1: []int{4, 5, 7},
			a2: []int{2, 5, 9},
			r:  []int{2, 4, 5, 5, 7, 9},
		},
		{
			a1: []int{3, 4, 6},
			a2: []int{1, 4, 5, 6, 9},
			r:  []int{1, 3, 4, 4, 5, 6, 6, 9},
		},
	}

	for _, v := range vs {
		r := Merge(v.a1, v.a2)
		s.Equal(v.r, r)
	}
}

func TestMergeTestSuite(t *testing.T) {
	s := &mergeTestSuite{}
	suite.Run(t, s)
}
