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
	"reflect"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type commonProviderTestSuite struct {
	suite.Suite
}

type mockStorageObject struct {
	mock.Mock
}

func (m *mockStorageObject) List() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}

func (m *mockStorageObject) Obtain(tag string) (*TagItem, error) {
	args := m.Called(tag)
	return args.Get(0).(*TagItem), args.Error(1)
}

func (s *commonProviderTestSuite) TestNewProviderOk() {
	tagItems := []*TagItem{
		{
			Tag:         "1",
			Max:         10,
			Min:         1,
			Description: "1",
		},
		{
			Tag:         "2",
			Max:         20,
			Min:         2,
			Description: "2",
		},
	}
	tagNames := make([]string, 0, len(tagItems))
	for _, item := range tagItems {
		tagNames = append(tagNames, item.Tag)
	}

	mockStorager := &mockStorageObject{}
	mockStorager.On("List").Return(tagNames, nil)
	for _, item := range tagItems {
		mockStorager.On("Obtain", item.Tag).Return(item, nil)
	}

	providerName := "test"
	p, err := NewProvider(providerName, mockStorager)
	s.NoError(err)
	s.NotNil(p)

	s.Equal(p.Name(), providerName)
	s.Equal(len(tagItems), len(p.tags))

	for _, item := range tagItems {
		v := p.tags[item.Tag]
		if !reflect.DeepEqual(item, v) {
			s.Failf("item equal failed", "expect %v, got %v", item, v)
		}
	}
}

func TestCommonProviderTestSuite(t *testing.T) {
	s := &commonProviderTestSuite{}
	suite.Run(t, s)
}
