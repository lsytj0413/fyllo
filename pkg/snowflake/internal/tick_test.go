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

type tickTestSuite struct {
	suite.Suite
}

func testNano(v uint64) func() uint64 {
	return func() uint64 {
		return v
	}
}

func (s *tickTestSuite) SetupSuite() {
	nano = testNano(Nano2MicroRatio)
}

func (s *tickTestSuite) TearDownSuite() {
	nano = timeNano
}

func (s *tickTestSuite) TestCurrent() {
	n := []uint64{Nano2MicroRatio, Nano2MicroRatio * 2}
	for i, v := range n {
		nano = testNano(v)
		s.Equal(uint64(i+1), defaultTicker.Current())
	}
}

func (s *tickTestSuite) TestWaitForNext() {
	n := []uint64{Nano2MicroRatio, Nano2MicroRatio * 2}
	for i, v := range n {
		nano = testNano(v)
		s.Equal(uint64(i+1), defaultTicker.WaitForNext(uint64(i)))
	}
}

func TestTickTestSuite(t *testing.T) {
	s := &tickTestSuite{}
	suite.Run(t, s)
}
