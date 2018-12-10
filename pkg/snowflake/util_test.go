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

package snowflake

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type utilTestSuite struct {
	suite.Suite
}

func (s *utilTestSuite) TestMaxMachineValue() {
	s.Equal(uint64(16), MaxMachineValue)
}

func (s *utilTestSuite) TestIsValidMachine() {
	for mid := uint64(0); mid < MaxMachineValue; mid++ {
		s.Nil(IsValidMachine(mid))
	}

	for mid, i := MaxMachineValue, 0; i < 10; i++ {
		s.Error(IsValidMachine(mid))
		mid += uint64(i)
		mid += uint64(i) * uint64(i)
	}
}

func (s *utilTestSuite) TestMaxTagValue() {
	s.Equal(uint64(256), MaxTagValue)
}

func (s *utilTestSuite) TestIsValidTag() {
	for bid := uint64(0); bid < MaxTagValue; bid++ {
		s.Nil(IsValidTag(bid))
	}

	for bid, i := MaxTagValue, 0; i < 10; i++ {
		s.Error(IsValidTag(bid))
		bid += uint64(i)
		bid += uint64(i) * uint64(i)
	}
}

func (s *utilTestSuite) TestMakeSnowflakeID() {
	type testCase struct {
		Description string

		timestamp      uint64
		mid            uint64
		bid            uint64
		sequenceNumber uint64

		expected uint64
	}
	testCases := []testCase{
		{
			Description:    "zero",
			timestamp:      StandardTimestamp,
			mid:            0,
			bid:            0,
			sequenceNumber: 0,
			expected:       0,
		},
		{
			Description:    "mid normal",
			timestamp:      StandardTimestamp,
			mid:            1,
			bid:            0,
			sequenceNumber: 0,
			expected:       0x0000000000040000,
		},
	}
	for _, tc := range testCases {
		s.Equal(tc.expected, MakeSnowflakeID(tc.timestamp, tc.mid, tc.bid, tc.sequenceNumber), tc.Description)
	}
}

func TestUtilTestSuite(t *testing.T) {
	s := &utilTestSuite{}
	suite.Run(t, s)
}
