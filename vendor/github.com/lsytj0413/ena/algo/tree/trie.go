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

// Trie is [Trie](https://en.wikipedia.org/wiki/Trie) or PrefixTree define
type Trie interface {
	Add(value string)
	Search(value string) bool
	StartsWith(value string) bool
	Size() int
}

type trie struct {
	root  *trieNode
	count int
}

func (p *trie) Size() int {
	return p.root.count
}

func (p *trie) Add(value string) {
	p.count++
	p.root.Add(value)
}

func (p *trie) Search(value string) bool {
	return p.root.Search(value)
}

func (p *trie) StartsWith(value string) bool {
	return p.root.StartsWith(value)
}

// NewTrie returns Trie implement
func NewTrie() Trie {
	return &trie{
		root: &trieNode{
			key:      "",
			count:    0,
			children: []*trieNode{},
		},
		count: 0,
	}
}
