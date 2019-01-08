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

package tree

import (
	"fmt"
	"strings"
)

type trieNode struct {
	key   string
	count int

	children []*trieNode
}

// CommonPrefix returns s1 and s2's equal prefix string
func CommonPrefix(s1 string, s2 string) string {
	// Keep s1.len <= s2.len
	if len(s1) > len(s2) {
		s1, s2 = s2, s1
	}

	for i := 0; i < len(s1); i++ {
		if s1[i] != s2[i] {
			return s1[0:i]
		}
	}
	return s1
}

// TrimPrefix delete prefix in s1, is arg prefix fit s1, it will return true
func TrimPrefix(s1 string, prefix string) (string, bool) {
	if len(prefix) <= len(s1) && s1[0:len(prefix)] == prefix {
		return s1[len(prefix):], true
	}

	return s1, false
}

func (n *trieNode) Add(value string) {
	if n.key == value {
		n.key = value
		n.count++
		return
	}

	if n.key == "" { // root node
		for _, child := range n.children {
			if child.key[0] == value[0] {
				child.Add(value)
				return
			}
		}

		n.children = append(n.children, &trieNode{
			key:      value,
			count:    1,
			children: []*trieNode{},
		})
		return
	}

	prefix := CommonPrefix(value, n.key)
	if prefix == "" { // no common prefix, panic
		panic(fmt.Sprintf("TrieNode Add: key %s and value %s has no common prefix", n.key, value))
	}

	switch {
	case prefix == value: // split current node
		n.children = []*trieNode{
			{
				key:      n.key[len(prefix):],
				count:    n.count,
				children: n.children,
			},
		}
		n.key = prefix
		n.count = 1
		return
	case prefix == n.key:
		value = value[len(prefix):]
		for _, child := range n.children {
			if child.key[0] == value[0] {
				child.Add(value)
				return
			}
		}
		n.children = append(n.children, &trieNode{
			key:      value,
			count:    1,
			children: []*trieNode{},
		})
		return
	}

	// split current node and has two child
	key := n.key[len(prefix):]
	value = value[len(prefix):]
	n.children = []*trieNode{
		{
			key:      key,
			count:    n.count,
			children: n.children,
		},
		{
			key:      value,
			count:    1,
			children: []*trieNode{},
		},
	}
	n.count = 0
	n.key = prefix
	return
}

func (n *trieNode) Search(value string) bool {
	if value == n.key {
		return n.count > 0
	}

	// len(value) > len(n.key), need search in child node if n.key is the prefix of value
	isPrefix := false
	if value, isPrefix = TrimPrefix(value, n.key); isPrefix {
		for _, child := range n.children {
			if child.key[0] == value[0] {
				return child.Search(value)
			}
		}
	}

	return false
}

func (n *trieNode) StartsWith(value string) bool {
	if len(value) == len(n.key) {
		return n.key == value
	}

	switch {
	case len(value) == len(n.key):
		return n.key == value
	case len(value) < len(n.key):
		return strings.HasPrefix(n.key, value)
	}

	// len(value) > len(n.key), process into child, n.key should be value's prefix
	isPrefix := false
	if value, isPrefix = TrimPrefix(value, n.key); isPrefix {
		for _, child := range n.children {
			if child.key[0] == value[0] {
				return child.StartsWith(value)
			}
		}
	}
	return false
}
