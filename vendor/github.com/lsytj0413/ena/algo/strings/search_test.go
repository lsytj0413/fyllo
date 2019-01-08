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

package strings

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type searchTestSuite struct {
	suite.Suite
}

func (s *searchTestSuite) TestFind() {
	cases := []struct {
		arg1   string
		arg2   string
		target int
	}{
		{
			arg1:   "abcd",
			arg2:   "ab",
			target: 0,
		},
		{
			arg1:   "abcd",
			arg2:   "cd",
			target: 2,
		},
		{
			arg1:   "abcd",
			arg2:   "ae",
			target: -1,
		},
		{
			arg1:   "abcd",
			arg2:   "ee",
			target: -1,
		},
	}

	for _, tc := range cases {
		s.Equal(tc.target, Find(tc.arg1, tc.arg2))
	}
}

func TestSearchTestSuite(t *testing.T) {
	s := &searchTestSuite{}
	suite.Run(t, s)
}
