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

package random

import (
	"context"

	"github.com/lsytj0413/fyllo/conf"
	ierror "github.com/lsytj0413/fyllo/error"
)

type createFunc func(conf.GeneratorConfig) (Generator, error)

// factory is design-pattern for create generator
type factory interface {
	register(createFunc)
	create(conf.GeneratorConfig) (Generator, error)
}

// Generator is interface for random generator
type Generator interface {
	Next(context.Context) (*conf.RandomResult, error)
}

type defFactory struct {
	fns []createFunc
}

func (d *defFactory) register(fn createFunc) {
	d.fns = append(d.fns, fn)
}

func (d *defFactory) create(config conf.GeneratorConfig) (Generator, error) {
	for _, fn := range d.fns {
		g, err := fn(config)

		if err == nil {
			return g, err
		}

		if !ierror.IsPluginNotImplement(err) {
			return nil, err
		}
	}

	return nil, ierror.NewPluginNotImplement("random factory doesn't have plugin")
}

var (
	f = &defFactory{
		fns: make([]createFunc, 0),
	}
)

// New will create a Generator or raise an error
func New(config conf.GeneratorConfig) (Generator, error) {
	return f.create(config)
}
