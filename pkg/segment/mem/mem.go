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

package mem

import (
	"fmt"

	"github.com/lsytj0413/fyllo/pkg/segment"
	"github.com/lsytj0413/fyllo/pkg/segment/internal"
)

const (
	// ProviderName for the mem segment provider
	ProviderName = "mem"
)

// memStroage implement internal.Storager, it storage item information in memory
type memStorage struct {
	tags map[string]*memRow
}

type memRow struct {
	Tag         string
	Max         uint64
	Step        uint64
	Description string
}

func (m *memStorage) List() ([]string, error) {
	tagNames := make([]string, 0, len(m.tags))

	for k := range m.tags {
		tagNames = append(tagNames, k)
	}

	return tagNames, nil
}

func (m *memStorage) Obtain(tag string) (*internal.TagItem, error) {
	rowItem, ok := m.tags[tag]
	if !ok {
		return nil, fmt.Errorf("not found")
	}

	r := &internal.TagItem{
		Tag:         tag,
		Min:         rowItem.Max,
		Max:         rowItem.Max + rowItem.Step,
		Description: rowItem.Description,
	}
	rowItem.Max = rowItem.Max + rowItem.Step + 1
	return r, nil
}

// Options is mem segment provider option
type Options struct {
	Args string
}

// NewProvider return mem segment provider implement
func NewProvider(options *Options) (segment.Provider, error) {
	return internal.NewProvider(ProviderName, &memStorage{})
}
