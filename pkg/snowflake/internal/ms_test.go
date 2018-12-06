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
	"time"

	"github.com/stretchr/testify/suite"
)

type msTestSuite struct {
	suite.Suite
}

func (s *msTestSuite) TestCurrMs() {
	now0 := uint64(time.Now().UnixNano()) / Nano2MicroRatio
	c := (&defMs{}).currMicroSeconds()
	now1 := uint64(time.Now().UnixNano()) / Nano2MicroRatio

	s.True(now0 <= c)
	s.True(now1 >= c)
}

func (s *msTestSuite) TestWaitMs() {
	now0 := uint64(time.Now().UnixNano()) / Nano2MicroRatio
	c := (&defMs{}).waitForNextMS(now0)
	now1 := uint64(time.Now().UnixNano()) / Nano2MicroRatio

	s.True(now0 < c)
	s.True(now1 >= c)
}

func TestMsTestSuite(t *testing.T) {
	p := &msTestSuite{}
	suite.Run(t, p)
}
