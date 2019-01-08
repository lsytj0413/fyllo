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
	"sort"
	"strconv"
	"sync/atomic"

	"github.com/lsytj0413/ena/algo/hash"
)

// Config is Consistent hash config
type Config struct {
	Replicas uint
}

const (
	// DefaultReplicas is default value for Config.Replicas
	DefaultReplicas = 32
)

// NewConfig returns Config object with Default Value
func NewConfig() *Config {
	return &Config{
		Replicas: DefaultReplicas,
	}
}

var (
	// ErrNoNodeValid is errors defines for Node Empty
	ErrNoNodeValid = errors.New("Empty Valid Node")
	// ErrDuplicateNodeKey is errors defines for Duplicated Node Key
	ErrDuplicateNodeKey = errors.New("Duplicated Node Key")
	// ErrNodeKeyNotExists is errors defines for Not Exists Node Key
	ErrNodeKeyNotExists = errors.New("Node Key Not Exists")
)

// Node is struct for backend
type Node struct {
	key    string
	weight uint
	Load   int32
	Extra  interface{}
}

// NewNode will returns a Node object
func NewNode(key string, weight uint, extra interface{}) *Node {
	return &Node{
		key:    key,
		weight: weight,
		Extra:  extra,
		Load:   0,
	}
}

// IncrLoad will increment Node.Load value
func (n *Node) IncrLoad(key string, delta int32) int32 {
	return atomic.AddInt32(&(n.Load), delta)
}

// Key return the Key value of node
func (n *Node) Key() string {
	return n.key
}

// Weight return the weight value of node
func (n *Node) Weight() uint {
	return n.weight
}

type virtualNode struct {
	key   string
	index uint32
	node  *Node
}

type ring []*virtualNode

// Len returns the length of the ring array.
func (x ring) Len() int { return len(x) }

// Less returns true if element i is less than element j.
func (x ring) Less(i, j int) bool { return x[i].index < x[j].index }

// Swap exchanges elements i and j.
func (x ring) Swap(i, j int) { x[i], x[j] = x[j], x[i] }

// Hasher is interface define for Consistent Hash
type Hasher interface {
	AddNode(n *Node) error
	RemoveNode(key string) error
	GetNode(key string) (*Node, error)
	Nodes() []*Node
	Hash(value string) (*Node, error)
}

type consistent struct {
	nodes  map[string]*Node
	circle ring

	c      *Config
	hasher hash.Hasher
}

// NewHasher will return Hasher object
func NewHasher(c *Config) Hasher {
	return &consistent{
		nodes:  make(map[string]*Node, 10),
		c:      c,
		circle: make([]*virtualNode, 0),
	}
}

func (h *consistent) AddNode(n *Node) error {
	_, exists := h.nodes[n.key]
	if exists {
		return ErrDuplicateNodeKey
	}

	h.nodes[n.key] = n

	for i := uint(0); i < n.weight*h.c.Replicas; i++ {
		key := h.eltKey(n.key, i)
		v := &virtualNode{
			key:   key,
			index: h.hasher.Uint32([]byte(key)),
			node:  n,
		}
		h.circle = append(h.circle, v)
	}
	sort.Sort(h.circle)

	return nil
}

func (h *consistent) eltKey(key string, index uint) string {
	return key + "#" + strconv.Itoa(int(index))
}

func (h *consistent) RemoveNode(key string) error {
	_, exists := h.nodes[key]
	if !exists {
		return ErrNodeKeyNotExists
	}

	delete(h.nodes, key)

	circle := make([]*virtualNode, 0)
	for _, v := range h.circle {
		if v.node.key != key {
			circle = append(circle, v)
		}
	}
	h.circle = circle

	return nil
}

func (h *consistent) GetNode(key string) (*Node, error) {
	n, exists := h.nodes[key]
	if !exists {
		return nil, ErrNodeKeyNotExists
	}

	return n, nil
}

func (h *consistent) Nodes() []*Node {
	r := make([]*Node, len(h.nodes))
	i := 0
	for _, v := range h.nodes {
		r[i] = v
		i++
	}
	return r
}

func (h *consistent) Hash(value string) (*Node, error) {
	if len(h.nodes) == 0 {
		return nil, ErrNoNodeValid
	}

	hashValue := h.hasher.Uint32([]byte(value))
	i := sort.Search(len(h.circle), func(x int) bool {
		return h.circle[x].index >= hashValue
	})
	if i >= len(h.circle) {
		i = 0
	}
	return h.circle[i].node, nil
}
