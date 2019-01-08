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

package common

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type arrayTest struct {
	suite.Suite
}

func (s *arrayTest) TestArray() {
	arr := [2]int{1, 2}
	err := ReverseSlice(&arr)
	s.NoError(err)

	s.Equal([2]int{2, 1}, arr)
}

func (s *arrayTest) TestSlice() {
	arr := []int{1, 2}
	err := ReverseSlice(&arr)
	s.NoError(err)

	s.Equal([]int{2, 1}, arr)
}

func TestArrayTestTestSuite(t *testing.T) {
	p := &arrayTest{}
	suite.Run(t, p)
}
