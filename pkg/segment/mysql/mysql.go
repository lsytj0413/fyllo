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

package mysql

import (
	"github.com/lsytj0413/fyllo/pkg/segment"
)

const (
	// ProviderName for the mysql segment provider
	ProviderName = "mysql"
)

type mysqlProvider struct {
}

func (p *mysqlProvider) Name() string {
	return ProviderName
}

func (p *mysqlProvider) Next(arg *segment.Arguments) (*segment.Result, error) {
	r := &segment.Result{
		Name: ProviderName,
		Next: 0,
		Labels: map[string]string{
			"tag": "",
		},
	}
	return r, nil
}

// Options is mysql segment provider option
type Options struct {
	Args string
}

// NewProvider return mysql segment provider implement
func NewProvider(options *Options) (segment.Provider, error) {
	return nil, nil
}
