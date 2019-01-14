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
	"github.com/stretchr/testify/mock"
)

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

func (m *mockStorageObject) ReverseCall() {
	for i, j := 0, len(m.ExpectedCalls)-1; i < j; i, j = i+1, j-1 {
		m.ExpectedCalls[i], m.ExpectedCalls[j] = m.ExpectedCalls[j], m.ExpectedCalls[i]
	}
	for i, j := 0, len(m.Calls)-1; i < j; i, j = i+1, j-1 {
		m.Calls[i], m.Calls[j] = m.Calls[j], m.Calls[i]
	}
}

func newMockStorager(tagItems []*TagItem) *mockStorageObject {
	tagNames := make([]string, 0, len(tagItems))
	for _, item := range tagItems {
		tagNames = append(tagNames, item.Tag)
	}

	mockStorager := &mockStorageObject{}
	mockStorager.On("List").Return(tagNames, nil)
	for _, item := range tagItems {
		mockStorager.On("Obtain", item.Tag).Return(item, nil)
	}
	return mockStorager
}
