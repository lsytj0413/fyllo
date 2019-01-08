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

package convert

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

type utilTestSuite struct {
	suite.Suite
}

func (s *utilTestSuite) TestPushOk() {
	tokens, err := ParseToLayoutField("%%123%acb%%c")

	fmt.Println(tokens)
	fmt.Println(err)
}

func TestUtilTestSuite(t *testing.T) {
	s := &utilTestSuite{}
	suite.Run(t, s)
}
