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
	"sync"
)

type safeConsistent struct {
	Hasher
	sync.Mutex
}

// NewSafeHasher will return Goroutine Safe Hasher object
func NewSafeHasher(c *Config) Hasher {
	return &safeConsistent{
		Hasher: NewHasher(c),
	}
}

func (h *safeConsistent) AddNode(n *Node) error {
	h.Lock()
	defer h.Unlock()

	return h.Hasher.AddNode(n)
}

func (h *safeConsistent) RemoveNode(key string) error {
	h.Lock()
	defer h.Unlock()

	return h.Hasher.RemoveNode(key)
}

func (h *safeConsistent) GetNode(key string) (*Node, error) {
	h.Lock()
	defer h.Unlock()

	return h.Hasher.GetNode(key)
}

func (h *safeConsistent) Nodes() []*Node {
	h.Lock()
	defer h.Unlock()

	return h.Hasher.Nodes()
}

func (h *safeConsistent) Hash(value string) (*Node, error) {
	h.Lock()
	defer h.Unlock()

	return h.Hasher.Hash(value)
}
