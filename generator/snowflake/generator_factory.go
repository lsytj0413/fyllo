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

package snowflake

import (
	"github.com/lsytj0413/fyllo/conf"
	"github.com/lsytj0413/fyllo/factory"
)

type snowflakeFactory struct {
}

func (f *snowflakeFactory) Name() string {
	return "UUIDRandomFactory"
}

func (f *snowflakeFactory) New(snowflakeConfig conf.GeneratorConfig) (conf.SnowflakeGenerator, error) {
	return &generatorMapper{}, nil
}

func init() {
	f := &snowflakeFactory{}
	err := factory.AppendFactory(f)
	if err != nil {
		panic(err)
	}
}
