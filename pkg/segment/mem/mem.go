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

	"github.com/lsytj0413/fyllo/pkg/common"
	"github.com/lsytj0413/fyllo/pkg/errors"
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
		return nil, errors.NewError(errors.EcodeSegmentQueryFailed, fmt.Sprintf("%s doesn't exists", tag))
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
		return nil, errors.NewError(errors.EcodeInitFailed, "segment provider[mem] argument should not be empty")
	}

	r := make(map[string]*memRow)
	argArray := strings.Split(args, ";")
	for _, arg := range argArray {
		row, err := parseProviderArgRow(arg)
		if err != nil {
			return nil, err
		}
		if _, ok := r[row.Tag]; ok {
			return nil, errors.NewError(errors.EcodeInitFailed, fmt.Sprintf("segment provider[mem] duplicate field %v", row.Tag))
		}
		r[row.Tag] = row
	}
	return r, nil
}

func parseProviderArgRow(arg string) (*memRow, error) {
	kvs, err := common.SplitKeyValueArrayString(arg)
	if err != nil {
		return nil, errors.NewError(errors.EcodeInitFailed, fmt.Sprintf("segment provider[mem] argument parse failed, %v", err))
	}

	mustFields := []string{"tag", "step"}
	for _, fieldName := range mustFields {
		if _, ok := kvs[fieldName]; !ok {
			return nil, errors.NewError(errors.EcodeInitFailed, fmt.Sprintf("segment provider[mem] argument parse failed, %s field not exists", fieldName))
		}
	}

	row := &memRow{}
	for key, value := range kvs {
		switch key {
		case "tag":
			row.Tag = value
		case "max":
			v, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return nil, errors.NewError(errors.EcodeInitFailed, fmt.Sprintf("segment provider[mem] argument parse failed, key[%v] value[%s] %v", key, value, err))
			}
			row.Max = v
		case "step":
			v, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return nil, errors.NewError(errors.EcodeInitFailed, fmt.Sprintf("segment provider[mem] argument parse failed, key[%v] value[%s] %v", key, value, err))
			}
			if v < 1 {
				return nil, errors.NewError(errors.EcodeInitFailed, fmt.Sprintf("segment provider[mem] argument parse failed, key[%v] value[%s] should bigger than 0", key, value))
			}
			row.Step = v
		case "desc":
			row.Description = value
		default:
			return nil, errors.NewError(errors.EcodeInitFailed, fmt.Sprintf("segment provider[mem] argument parse failed, key[%v] is invalid", key))
		}
	}

	return row, nil
}
