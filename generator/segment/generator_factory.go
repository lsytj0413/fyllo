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

package segment

import (
	"github.com/lsytj0413/fyllo/conf"
	"github.com/lsytj0413/fyllo/factory"
)

type segmentFactory struct {
}

func (f *segmentFactory) Name() string {
	return "UUIDRandomFactory"
}

func (f *segmentFactory) New(segmentConfig conf.GeneratorConfig) (conf.SegmentGenerator, error) {
	return &generator{}, nil
}

func init() {
	f := &segmentFactory{}
	err := factory.AppendFactory(f)
	if err != nil {
		panic(err)
	}
}
