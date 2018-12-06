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

package fs

import "fmt"

// Node is the external representation of the inode with additional fields
type Node struct {
	Key   string
	Value *string
	Dir   bool
	Nodes NodeArray
}

// String implements fmt.Stringer
func (n Node) String() string {
	var value string
	if n.Value != nil {
		value = *n.Value
	}

	return fmt.Sprintf("Node(Key=%s, Value=%s, Dir=%v)", n.Key, value, n.Dir)
}

// Clone return node clone object
// If the node is a directory, it will clone all the content under this directory
// If the node is a file, it will clone the file
func (n *Node) Clone() *Node {
	if n == nil {
		return nil
	}

	nn := &Node{
		Key: n.Key,
		Dir: n.Dir,
	}

	if n.Value != nil {
		nn.Value = &(*n.Value)
	}

	if n.Nodes != nil {
		nn.Nodes = nn.Nodes.Clone()
	}

	return nn
}

// NodeArray is list of Node
type NodeArray []*Node

// Len implements sort.Interface
func (na NodeArray) Len() int {
	return len(na)
}

// Less implements sort.Interface
func (na NodeArray) Less(i int, j int) bool {
	return na[i].Key < na[j].Key
}

// Swap implements sort.Interface
func (na NodeArray) Swap(i int, j int) {
	na[i], na[j] = na[j], na[i]
}

// Clone object
func (na NodeArray) Clone() NodeArray {
	nodes := make(NodeArray, na.Len())
	for i, n := range na {
		nodes[i] = n.Clone()
	}

	return nodes
}
