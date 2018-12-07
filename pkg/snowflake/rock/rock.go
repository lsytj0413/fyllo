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

// Package rock implement the provider which has const machine id
package rock

import (
	"github.com/lsytj0413/fyllo/pkg/snowflake"
)

const (
	// ProviderName for the const snowflake provider
	ProviderName = "rock"
)

type rockProvider struct {
}

func (p *rockProvider) Name() string {
	return ProviderName
}

func (p *rockProvider) Next() (uint64, error) {
	return 0, nil
}

// Options is rock snowflake provider option
type Options struct {
	Args string
}

// NewProvider return rock snowflake provider implement
func NewProvider(options *Options) (snowflake.Provider, error) {
	return nil, nil
}
