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

package builder

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/lsytj0413/fyllo/pkg/segment"
	"github.com/lsytj0413/fyllo/pkg/segment/internal"
	"github.com/lsytj0413/fyllo/pkg/segment/mem"
	"github.com/lsytj0413/fyllo/pkg/segment/mysql"
)

type builderTestSuite struct {
	suite.Suite
}

func (s *builderTestSuite) TestNewBuilderOk() {
	names := []string{mem.ProviderName, mysql.ProviderName}
	for _, name := range names {
		b, err := NewBuilder(&Options{
			ProviderName: name,
		})
		s.NoError(err)
		s.NotNil(b)
	}
}
func (s *builderTestSuite) TestNewBuilderFailed() {
	names := []string{"test"}
	for _, name := range names {
		b, err := NewBuilder(&Options{
			ProviderName: name,
		})
		s.Error(err)
		s.Nil(b)
	}
}

func (s *builderTestSuite) TestBuildOk() {
	options := []*Options{
		{
			ProviderName: mem.ProviderName,
			ProviderArgs: "tag=1,step=100",
		},
		{
			ProviderName: mysql.ProviderName,
			ProviderArgs: "",
		},
	}
	oldCreateMysqlProvider, oldCreateMemProvider := createMysqlProvider, createMemProvider
	createMysqlProvider = func(*Options) (segment.Provider, error) {
		return &internal.CommonProvider{}, nil
	}
	createMemProvider = func(*Options) (segment.Provider, error) {
		return &internal.CommonProvider{}, nil
	}
	for _, option := range options {
		b, err := NewBuilder(option)
		s.NoError(err)

		p, err := b.Build()
		s.NoError(err)
		s.NotNil(p)
	}
	createMysqlProvider, createMemProvider = oldCreateMysqlProvider, oldCreateMemProvider
}

func (s *builderTestSuite) TestBuildFailed() {
	options := []*Options{
		{
			ProviderName: mem.ProviderName,
			ProviderArgs: "",
		},
	}
	for _, option := range options {
		b, err := NewBuilder(option)
		s.NoError(err)

		option.ProviderName = "test"

		p, err := b.Build()
		s.Error(err)
		s.Nil(p)
	}
}

func TestBuilderTestSuite(t *testing.T) {
	s := &builderTestSuite{}
	suite.Run(t, s)
}
