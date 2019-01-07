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
	"strconv"
	"strings"

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
	tags, err := parseProviderArgs(options.Args)
	if err != nil {
		return nil, err
	}

	return internal.NewProvider(ProviderName, &memStorage{
		tags: tags,
	})
}

func parseProviderArgs(args string) (map[string]*memRow, error) {
	if 0 == len(args) {
		return nil, fmt.Errorf("empty args")
	}

	r := make(map[string]*memRow)
	argArray := strings.Split(args, ";")
	for _, arg := range argArray {
		row, err := parseProviderArgRow(arg)
		if err != nil {
			return nil, err
		}
		if _, ok := r[row.Tag]; ok {
			return nil, fmt.Errorf("duplicate %s tag name", row.Tag)
		}
		r[row.Tag] = row
	}
	return r, nil
}

func parseProviderArgRow(arg string) (*memRow, error) {
	fields := strings.Split(arg, ",")

	valueMap := make(map[string]string)
	for _, field := range fields {
		arr := strings.Split(field, "=")
		if len(arr) != 2 {
			return nil, fmt.Errorf("wrong field %s", field)
		}

		key, value := strings.ToLower(arr[0]), arr[1]
		if _, ok := valueMap[key]; ok {
			return nil, fmt.Errorf("duplicate key %s", key)
		}
		valueMap[key] = value
	}

	mustFields := []string{"tag", "step"}
	for _, fieldName := range mustFields {
		if _, ok := valueMap[fieldName]; !ok {
			return nil, fmt.Errorf("%s key loss", fieldName)
		}
	}

	row := &memRow{}
	for key, value := range valueMap {
		switch key {
		case "tag":
			row.Tag = value
		case "max":
			v, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("%s key %s value parse failed, %v", key, value, err)
			}
			row.Max = v
		case "step":
			v, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("%s key %s value parse failed, %v", key, value, err)
			}
			if v < 1 {
				return nil, fmt.Errorf("%s key %s value should bigger than 0", key, value)
			}
			row.Step = v
		case "desc":
			row.Description = value
		default:
			return nil, fmt.Errorf("%s key wrong name", key)
		}
	}

	return row, nil
}
