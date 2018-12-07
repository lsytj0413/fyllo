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
	"errors"
	"strconv"

	"github.com/lsytj0413/fyllo/pkg/snowflake"
	"github.com/lsytj0413/fyllo/pkg/snowflake/internal"
)

const (
	// ProviderName for the const snowflake provider
	ProviderName = "rock"
)

type rockIdentifier struct {
	mid uint64
}

func (r *rockIdentifier) Identify() (uint64, error) {
	return r.mid, nil
}

// Options is rock snowflake provider option
type Options struct {
	Args string
}

// NewProvider return rock snowflake provider implement
func NewProvider(options *Options) (snowflake.Provider, error) {
	mid, err := strconv.Atoi(options.Args)
	if err != nil {
		return nil, err
	}
	if mid < 0 || uint64(mid) >= snowflake.MaxMachineValue {
		return nil, errors.New("wrong machine id")
	}

	return internal.NewProvider(ProviderName, &rockIdentifier{
		mid: uint64(mid),
	}), nil
}
