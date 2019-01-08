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
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type reflectEnvTest struct {
	suite.Suite
}

func (s *reflectEnvTest) TestString() {
	src := "${PATH}"
	err := ReplaceEnv(&src)
	s.NoError(err)
	s.Equal(os.Getenv("PATH"), src)
}

type St struct {
	Name string
}

func (s *reflectEnvTest) TestTopStruct() {
	src := struct {
		Name  string
		Value int
		Name2 string
		S     St
	}{
		Name:  "${PATH}",
		Value: 1,
		Name2: "${",
		S: St{
			Name: "${PATH}",
		},
	}
	err := ReplaceEnv(&src)
	s.NoError(err)
	s.Equal(os.Getenv("PATH"), src.Name)
	s.Equal("${", src.Name2)
	s.Equal(1, src.Value)
	s.Equal(os.Getenv("PATH"), src.S.Name)
}

func TestReflectEnvTestSuite(t *testing.T) {
	p := &reflectEnvTest{}
	suite.Run(t, p)
}
