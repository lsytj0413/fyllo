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
	"context"

	"github.com/lsytj0413/fyllo/conf"
	ierror "github.com/lsytj0413/fyllo/error"
	"github.com/lsytj0413/fyllo/snowflake/internal"
)

type constPlugin struct {
	g internal.Generator
}

func (p *constPlugin) Next(c context.Context, tag uint64) (*conf.SnowflakeResult, error) {
	r, err := p.g.Next(tag)
	if err != nil {
		return nil, err
	}

	return &conf.SnowflakeResult{
		Tag:       r.Tag,
		Machine:   r.Machine,
		Next:      r.Next,
		Timestamp: r.Timestamp,
		Sequence:  r.Sequence,
	}, nil
}

func init() {
	createFn := func(config conf.GeneratorConfig) (Generator, error) {
		if config.Plugin != "CONST" {
			return nil, ierror.NewPluginNotImplement("CONST plugin only for CONST machine")
		}

		g, err := internal.NewGenerator(1)
		if err != nil {
			return nil, err
		}
		return &constPlugin{g: g}, nil
	}

	f.register(createFn)
}
