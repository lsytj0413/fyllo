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

package logger

import (
	"fmt"
	"testing"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/stretchr/testify/suite"
)

type layoutFormatterTestSuite struct {
	suite.Suite
}

func (s *layoutFormatterTestSuite) TestPushOk() {
	f, err := NewLayoutFormatter("[%d] [%level] [%P:%M:%L:%F]^sfd %msg xs%%re")
	if err != nil {
		fmt.Println(err)
		s.Fail("err should be nil")
	}

	e := &logrus.Entry{
		Time:    time.Now(),
		Level:   logrus.Level(0),
		Message: "layout test",
	}
	l, err := f.Format(e)
	s.NoError(err)

	fmt.Println(string(l))
}

func TestLayoutFormatterTestSuite(t *testing.T) {
	s := &layoutFormatterTestSuite{}
	suite.Run(t, s)
}
