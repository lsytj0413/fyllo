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
	"bytes"
	"errors"
	"fmt"
)

// IsAlpha check byte is letter
func IsAlpha(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}

// ParseToLayoutField parse layout to LayoutFields
func ParseToLayoutField(layout string) ([]*LayoutField, error) {
	fields := make([]*LayoutField, 0)

	// TODO: merge same type layout field
	n, mode := bytes.Buffer{}, 0
	for i := 0; i < len(layout); i++ {
		switch mode {
		case 0: // init mode
			if layout[i] == '%' {
				mode = 1 // field key start
				if n.Len() != 0 {
					fields = append(fields, &LayoutField{
						Type:  LayoutFieldText,
						Value: n.String(),
					})
					n.Reset()
				}
			} else {
				n.WriteByte(layout[i])
			}
		case 1: // last char is %, converter key start
			switch {
			case layout[i] == '%':
				// %% pattern
				fields = append(fields, &LayoutField{
					Type:  LayoutFieldText,
					Value: "%",
				})
				mode = 0
			case IsAlpha(layout[i]):
				// converter key start
				n.WriteByte(layout[i])
				mode = 2
			default:
				return nil, fmt.Errorf("unexpected char[%v] after %% at %d", layout[i], i)
			}
		case 2:
			if IsAlpha(layout[i]) {
				// converter key continue
				n.WriteByte(layout[i])
			} else {
				mode = 0
				fields = append(fields, &LayoutField{
					Type:  LayoutFieldConverter,
					Value: n.String(),
				})
				n.Reset()
				i-- // re process this char
			}
		}
	}

	if mode == 1 {
		return nil, errors.New("unexpected % at end of pattern")
	}

	// last field
	if n.Len() != 0 {
		if mode == 0 {
			fields = append(fields, &LayoutField{
				Type:  LayoutFieldText,
				Value: n.String(),
			})
		} else if mode == 2 {
			fields = append(fields, &LayoutField{
				Type:  LayoutFieldConverter,
				Value: n.String(),
			})
		}
	}

	return fields, nil
}
