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
	"hash/fnv"

	"github.com/lsytj0413/fyllo/conf"
	uuid "github.com/satori/go.uuid"
)

type uuidGenerator struct {
}

func (g *uuidGenerator) Next(c context.Context) (*conf.RandomResult, error) {
	r := &conf.RandomResult{}
	r.Identify = uuid.NewV4().String()

	h := fnv.New64a()
	_, err := h.Write([]byte(r.Identify))
	if err != nil {
		return nil, err
	}
	r.Next = h.Sum64()
	return r, nil
}

func init() {
	createFn := func(config conf.GeneratorConfig) (Generator, error) {
		if config.Plugin != "UUID" {
			return nil, errNotImplement
		}

		return &uuidGenerator{}, nil
	}

	f.register(createFn)
}
