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

package rock

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/lsytj0413/fyllo/pkg/snowflake"
)

type rockTestSuite struct {
	suite.Suite
}

func (s *rockTestSuite) TestNewProviderOk() {
	provider, err := NewProvider(&Options{
		Args: "0",
	})
	s.NoError(err)
	s.NotNil(provider)
	s.Equal(ProviderName, provider.Name())
}

func (s *rockTestSuite) TestNewProviderArgsFail() {
	type testCase struct {
		Description string
		Args        string
	}
	testCases := []testCase{
		{
			Description: "Args Not Number",
			Args:        "a",
		},
		{
			Description: "Args negative",
			Args:        "-1",
		},
		{
			Description: "Args equal max",
			Args:        fmt.Sprintf("%d", snowflake.MaxMachineValue),
		},
		{
			Description: "Args bigger than max",
			Args:        fmt.Sprintf("%d", snowflake.MaxMachineValue+1),
		},
	}

	for _, tc := range testCases {
		provider, err := NewProvider(&Options{
			Args: tc.Args,
		})
		s.Error(err, tc.Description)
		s.Nil(provider, tc.Description)
	}
}

func (s *rockTestSuite) TestRockIdentifier() {
	for i := uint64(0); i < snowflake.MaxMachineValue; i += 3 {
		identifier := &rockIdentifier{
			mid: i,
		}
		actual, err := identifier.Identify()
		s.NoError(err)
		s.Equal(i, actual)
	}
}

func (s *rockTestSuite) TestProviderNextOk() {
	provider, _ := NewProvider(&Options{
		Args: "0",
	})

	for tag := uint64(0); tag < snowflake.MaxTagValue; tag++ {
		r, err := provider.Next(&snowflake.Arguments{
			Tag: tag,
		})
		s.NoError(err)
		s.True(r.Next > 0)
		s.Equal(ProviderName, r.Name)

		timestamp, err := strconv.ParseUint(r.Labels[snowflake.LabelTimestamp], 10, 64)
		s.NoError(err)
		sequenceNumber, err := strconv.ParseUint(r.Labels[snowflake.LabelSequence], 10, 64)
		s.NoError(err)
		rtag, err := strconv.ParseUint(r.Labels[snowflake.LabelTag], 10, 64)
		s.NoError(err)
		s.Equal(tag, rtag)
		machine, err := strconv.ParseUint(r.Labels[snowflake.LabelMachine], 10, 64)
		s.NoError(err)
		s.Equal(uint64(0), machine)

		s.Equal(snowflake.MakeSnowflakeID(timestamp, machine, rtag, sequenceNumber), r.Next)
		break
	}
}

func TestRockTestSuite(t *testing.T) {
	s := &rockTestSuite{}
	suite.Run(t, s)
}
