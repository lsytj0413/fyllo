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

package cons

import (
	"errors"

	"github.com/lsytj0413/ena/algo/hash"
)

// Carper is interface for [CARP](https://tools.ietf.org/html/draft-vinod-carp-v1-03) algorithm
type Carper interface {
	Hash(string) (string, error)
}

type carp struct {
	endpoints []string
	hasher    hash.Hasher
}

// NewCarp will returns Carper implement
func NewCarp(endpoints []string) (Carper, error) {
	if 0 == len(endpoints) {
		return nil, errors.New("endpoints length must be greater than 0")
	}

	return &carp{
		endpoints: endpoints,
		hasher:    hash.NewHash(),
	}, nil
}

func (h *carp) Hash(key string) (string, error) {
	if 1 == len(h.endpoints) {
		return h.endpoints[0], nil
	}

	hashedArr := make([]uint64, len(h.endpoints))
	for i, endpoint := range h.endpoints {
		hashedArr[i] = h.hasher.Uint64([]byte(key + endpoint))
	}

	min, endpoint := hashedArr[0], h.endpoints[0]
	for i, v := range hashedArr {
		if v < min {
			endpoint = h.endpoints[i]
			min = v
		}
	}

	return endpoint, nil
}
