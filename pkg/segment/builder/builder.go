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

package builder

import (
	"fmt"
	"strings"

	"github.com/lsytj0413/fyllo/pkg/errors"
	"github.com/lsytj0413/fyllo/pkg/segment"
	"github.com/lsytj0413/fyllo/pkg/segment/mem"
	"github.com/lsytj0413/fyllo/pkg/segment/mysql"
)

// AvailableProviders supported by the segment provider builder.
var AvailableProviders = []string{
	mysql.ProviderName,
	mem.ProviderName,
}

// AvailableProvidersDescription is string readable description for providers list
var AvailableProvidersDescription string

// Options for builder
type Options struct {
	ProviderName string
	ProviderArgs string
}

// Builder for build segment provider
type Builder interface {
	// Build will return segment.Provider implement specfied by the ProviderName
	Build() (segment.Provider, error)
}

// NewBuilder return Builder instance
func NewBuilder(options *Options) (Builder, error) {
	found := false
	for _, name := range AvailableProviders {
		if name == options.ProviderName {
			found = true
			break
		}
	}
	if !found {
		return nil, errors.NewError(errors.EcodeProviderNotImplement, fmt.Sprintf("Invalid Segment ProviderName[%s], Avaliable: %s", options.ProviderName, AvailableProvidersDescription))
	}

	return &builder{
		options: options,
	}, nil
}

type builder struct {
	options *Options
}

func (b *builder) Build() (segment.Provider, error) {
	switch b.options.ProviderName {
	case mysql.ProviderName:
		return createMysqlProvider(b.options)
	case mem.ProviderName:
		return createMemProvider(b.options)
	}

	return nil, errors.NewError(errors.EcodeProviderNotImplement, fmt.Sprintf("Invalid Segment ProviderName[%s], Avaliable: %s", b.options.ProviderName, AvailableProvidersDescription))
}

var createMysqlProvider = func(option *Options) (segment.Provider, error) {
	return mysql.NewProvider(&mysql.Options{
		Args: option.ProviderArgs,
	})
}

var createMemProvider = func(option *Options) (segment.Provider, error) {
	return mem.NewProvider(&mem.Options{
		Args: option.ProviderArgs,
	})
}

func init() {
	AvailableProvidersDescription = "[" + strings.Join(AvailableProviders, ", ") + "]"
}
