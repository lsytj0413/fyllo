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

package uuid

import (
	gouuid "github.com/satori/go.uuid"

	"github.com/lsytj0413/ena/algo/hash"
	"github.com/lsytj0413/fyllo/pkg/random"
)

const (
	// ProviderName for the uuid random provider
	ProviderName = "uuid"
)

type uuidProvider struct {
}

func (p *uuidProvider) Name() string {
	return ProviderName
}

func (p *uuidProvider) Next(arg *random.Arguments) (*random.Result, error) {
	identify := gouuid.NewV4().String()

	r := &random.Result{
		Name: ProviderName,
		Labels: map[string]string{
			random.LabelIdentify: identify,
		},
		Next: hash.Uint64([]byte(identify)),
	}
	return r, nil
}

// Options is uuid random provider option
type Options struct {
	Args string
}

// NewProvider return uuid random provider implement
func NewProvider(options *Options) (random.Provider, error) {
	return &uuidProvider{}, nil
}
