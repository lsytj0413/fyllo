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
	"fmt"
)

// Builder build layout strings to Converter
type Builder interface {
	// Build construct Converter from layout pattern string
	Build(layout string) (Converter, error)

	// AddFieldBuilder add user field converter builder, if key conflict it will overwrite old field converter builder
	AddFieldBuilder(f FieldConverterBuilder)
}

// FieldConverterBuilder build FieldConverer
type FieldConverterBuilder interface {
	Build(layoutField *LayoutField) (FieldConverter, error)
	Key() string
}

// defBuilderImpl implement Builder interface
type defBuilderImpl struct {
	fieldBuilers map[string]FieldConverterBuilder
}

// NewBuilder return Builder
func NewBuilder() Builder {
	return &defBuilderImpl{
		fieldBuilers: make(map[string]FieldConverterBuilder),
	}
}

func (b *defBuilderImpl) Build(layout string) (Converter, error) {
	fields, err := ParseToLayoutField(layout)
	if err != nil {
		return nil, err
	}

	fieldConverters := make([]FieldConverter, 0, len(fields))
	for _, field := range fields {
		switch {
		case field.Type == LayoutFieldText:
			fieldConverters = append(fieldConverters, &textConverter{
				text: field.Value,
			})
		case field.Type == LayoutFieldConverter:
			var converter FieldConverter
			if builder, ok := b.fieldBuilers[field.Value]; ok {
				converter, err = builder.Build(field)
				if err != nil {
					return nil, err
				}
				fieldConverters = append(fieldConverters, converter)
			} else {
				// TODO: builder pattern?
				switch field.Value {
				case string(FieldKeyDate):
					fieldConverters = append(fieldConverters, &dateConverter{
						timestampFormat: "2006-01-02T15:04:05.000000000Z07:00",
					})
				case string(FieldKeyLevel):
					fieldConverters = append(fieldConverters, &levelConverter{})
				case string(FieldKeyPackage):
					fieldConverters = append(fieldConverters, &packageConverter{})
				case string(FieldKeyFile):
					fieldConverters = append(fieldConverters, &fileConverter{})
				case string(FieldKeyMethod):
					fieldConverters = append(fieldConverters, &methodConverter{})
				case string(FieldKeyLine):
					fieldConverters = append(fieldConverters, &lineConverter{})
				case string(FieldKeyMessage):
					fieldConverters = append(fieldConverters, &messageConverter{})
				default:
					return nil, fmt.Errorf("unsupported field key: %s", field.Value)
				}
			}
		default:
			return nil, fmt.Errorf("unexpected field type %v", field.Type)
		}
	}

	return NewConverter(fieldConverters), nil
}

func (b *defBuilderImpl) AddFieldBuilder(f FieldConverterBuilder) {
	b.fieldBuilers[f.Key()] = f
}
