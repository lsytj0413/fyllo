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

package factory

import (
	"fmt"
	"reflect"

	"github.com/lsytj0413/fyllo/conf"
)

type factory interface {
	Name() string
}

// SnowflakeFactory is interface define for SnowflakeGenerator construct factory
type SnowflakeFactory interface {
	factory
	New() (conf.SnowflakeGenerator, error)
}

// SegmentFactory is interface define for SegementGenerator construct factory
type SegmentFactory interface {
	factory
	New() (conf.SegmentGenerator, error)
}

// RandomFactory is interface define for RandomGenerator construct factory
type RandomFactory interface {
	factory
	New() (conf.RandomGenerator, error)
}

var (
	snowflakeFactorys []SnowflakeFactory
	segmentFactorys   []SegmentFactory
	randomFactorys    []RandomFactory
)

// AppendFactory will add factory instance
func AppendFactory(fac factory) error {
	switch fac := fac.(type) {
	case SnowflakeFactory:
		return AppendSnowflakeFactory(fac)
	case SegmentFactory:
		return AppendSegmentFactory(fac)
	case RandomFactory:
		return AppendRandomFactory(fac)
	default:
		return fmt.Errorf("factory.AppendFactory %v: unsupport type", reflect.TypeOf(fac).String())
	}
}

// AppendSnowflakeFactory will add SnowflakeFactory instance
func AppendSnowflakeFactory(fac SnowflakeFactory) error {
	for _, f := range snowflakeFactorys {
		if f.Name() == fac.Name() {
			return fmt.Errorf("factory.AppendSnowflakeFactory %v, Duplicate Name", f.Name())
		}
	}

	snowflakeFactorys = append(snowflakeFactorys, fac)
	return nil
}

// AppendSegmentFactory will add SegmentFactory instance
func AppendSegmentFactory(fac SegmentFactory) error {
	for _, f := range segmentFactorys {
		if f.Name() == fac.Name() {
			return fmt.Errorf("factory.AppendSegmentFactory %v, Duplicate Name", f.Name())
		}
	}

	segmentFactorys = append(segmentFactorys, fac)
	return nil
}

// AppendRandomFactory will add RandomFactory instance
func AppendRandomFactory(fac RandomFactory) error {
	for _, f := range randomFactorys {
		if f.Name() == fac.Name() {
			return fmt.Errorf("factory.AppendRandomFactory %v, Duplicate Name", f.Name())
		}
	}

	randomFactorys = append(randomFactorys, fac)
	return nil
}
