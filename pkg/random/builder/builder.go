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

	ierror "github.com/lsytj0413/fyllo/pkg/error"
	"github.com/lsytj0413/fyllo/pkg/random"
	"github.com/lsytj0413/fyllo/pkg/random/uuid"
)

// AvailableProviders supported by the random provider builder.
var AvailableProviders = []string{
	uuid.ProviderName,
}

// AvailableProvidersDescription is string readable description for providers list
var AvailableProvidersDescription string

// Options for builder
type Options struct {
	ProviderName string
	ProviderArgs string
}

// Builder for build random provider
type Builder interface {
	// Build will return random.Provider implement specfied by the ProviderName
	Build() (random.Provider, error)
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
		return nil, ierror.NewError(ierror.EcodeProviderNotImplement, fmt.Sprintf("Invalid Random ProviderName[%s], Avaliable: %s", options.ProviderName, AvailableProvidersDescription))
	}

	return &builder{
		options: options,
	}, nil
}

type builder struct {
	options *Options
}

func (b *builder) Build() (random.Provider, error) {
	switch b.options.ProviderName {
	case uuid.ProviderName:
		return uuid.NewProvider(&uuid.Options{
			Args: b.options.ProviderArgs,
		})
	}

	return nil, ierror.NewError(ierror.EcodeProviderNotImplement, fmt.Sprintf("Invalid Random ProviderName[%s], Avaliable: %s", b.options.ProviderName, AvailableProvidersDescription))
}

func init() {
	AvailableProvidersDescription = "[" + strings.Join(AvailableProviders, ", ") + "]"
}
