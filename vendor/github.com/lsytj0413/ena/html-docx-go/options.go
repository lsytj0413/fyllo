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

package html2docx

import "errors"

// Margins defines Docx margin option
type Margins struct {
	Top    uint32
	Right  uint32
	Bottom uint32
	Left   uint32
	Header uint32
	Footer uint32
	Gutter uint32
}

const (
	// Landscape is orient landscape value
	Landscape = "landscape"
	// Portrait is orient portrait value
	Portrait = "portrait"
)

// Option defines Docx optional
type Option struct {
	Margin Margins
	Orient string
	Width  uint32
	Height uint32
}

var (
	defaultOption = Option{
		Margin: Margins{
			Top:    1440,
			Right:  1440,
			Bottom: 1440,
			Left:   1440,
			Header: 720,
			Footer: 720,
			Gutter: 0,
		},
		Orient: Portrait,
		Width:  12240,
		Height: 15840,
	}

	// ErrInvalidOrientValue is errors for invalid orient value
	ErrInvalidOrientValue = errors.New("Invalid Orient Value")
)

// NewOption will returns a Option instance with default option
func NewOption() (Option, error) {
	return defaultOption, nil
}

// NewOptionFromOrient will returns a Option instance with given orient and default option
func NewOptionFromOrient(orient string) (Option, error) {
	switch orient {
	case Portrait:
		return NewOption()
	case Landscape:
		option, err := NewOption()
		if err != nil {
			return defaultOption, err
		}
		option.Orient = orient
		option.Width = 15840
		option.Height = 12240
		return option, nil
	}

	return defaultOption, ErrInvalidOrientValue
}
