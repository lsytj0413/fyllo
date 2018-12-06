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

import (
	"fmt"
	"time"
)

// inode is basic element in the store system
type inode struct {
	Path  string
	Value string

	Parent   *inode
	Children map[string]*inode // for directory

	Expire time.Time
	// A reference to the store this inode is attached to
	store *defFileSystemStore
}

// String implements fmt.Stringer
func (n inode) String() string {
	return fmt.Sprintf("inode(Path=%s, Value=%s)", n.Path, n.Value)
}

func newFileInode(store *defFileSystemStore, nodePath string, value string, parent *inode) *inode {
	return &inode{
		Path:   nodePath,
		Value:  value,
		Parent: parent,
		store:  store,
	}
}

func newDirInode(store *defFileSystemStore, nodePath string, parent *inode) *inode {
	return &inode{
		Path:     nodePath,
		Parent:   parent,
		Children: make(map[string]*inode),
		store:    store,
	}
}

// IsHidden check the inode is hidden, which inode name start with dot is a hidden inode
func (n *inode) IsHidden() bool {
	v := name(n.Path)
	if v == "" {
		// TODO: maybe panic?
		return false
	}

	return v[0] == '.'
}

// IsDir check the inode is directory
func (n *inode) IsDir() bool {
	return n.Children != nil
}

// Read function gets the value of the node
// If node is a directory, EcodeNotFile return
func (n *inode) Read() (string, error) {
	if n.IsDir() {
		return "", NewError(EcodeNotFile, fmt.Sprintf("Read %s", n.Path))
	}

	return n.Value, nil
}

// Write function set the value of the node to the given value
// If node is a directory, EcodeNotFile return
func (n *inode) Write(value string) error {
	if n.IsDir() {
		return NewError(EcodeNotFile, fmt.Sprintf("Write %s", n.Path))
	}

	n.Value = value
	return nil
}

// List return child inode at the current inode
// If current inode is file, EcodeNotDir error return
func (n *inode) List() ([]*inode, error) {
	if !n.IsDir() {
		return nil, NewError(EcodeNotDir, fmt.Sprintf("List %s", n.Path))
	}

	nodes := make([]*inode, len(n.Children))

	i := 0
	for _, node := range n.Children {
		nodes[i] = node
		i++
	}

	return nodes, nil
}

// GetChild returns the child inode under the directory inode
// If current inode is file, EcodeNotDir return
// If child not exists, returns error
func (n *inode) GetChild(name string) (*inode, error) {
	if !n.IsDir() {
		return nil, NewError(EcodeNotDir, fmt.Sprintf("GetChild %s: child=%s", n.Path, name))
	}

	child, ok := n.Children[name]
	if ok {
		return child, nil
	}

	return nil, NewError(EcodeNotExists, fmt.Sprintf("GetChild %s: child=%s", n.Path, name))
}

// Add function adds a inode to the directory inode
// If current inode is not directory, returns EcodeNotDir
// If same name already exists under the directory, returns EcodeFileExists
func (n *inode) Add(child *inode) error {
	if !n.IsDir() {
		return NewError(EcodeNotDir, fmt.Sprintf("Add %s: child=%s", n.Path, child.Path))
	}

	name := name(key(child.Path))
	if _, ok := n.Children[name]; ok {
		return NewError(EcodeExists, fmt.Sprintf("Add %s: child=%s", n.Path, child.Path))
	}

	n.Children[name] = child
	return nil
}

// Remove function remove the node
// If current inode is directory and dir is false, return EcodeNotFile
// If current inode is directory and has child and recursive is false, return EcodeDirNotEmpty
func (n *inode) Remove(dir bool, recursive bool) error {
	if !n.IsDir() {
		name := name(key(n.Path))
		if n.Parent != nil && n.Parent.Children[name] == n {
			delete(n.Parent.Children, name)
		}

		return nil
	}

	if !dir {
		return NewError(EcodeNotFile, fmt.Sprintf("Remove %s, dir must be true", n.Path))
	}

	if len(n.Children) != 0 && !recursive {
		return NewError(EcodeDirNotEmpty, fmt.Sprintf("Remove %s, recursive must be true", n.Path))
	}

	for _, child := range n.Children {
		child.Remove(true, true)
	}

	// Delete self
	name := name(key(n.Path))
	if n.Parent != nil && n.Parent.Children[name] == n {
		delete(n.Parent.Children, name)
	}
	return nil
}

// Clone return inode clone object
// If the node is a directory, it will clone all the content under this directory
// If the node is a file, it will clone the file
func (n *inode) Clone() *inode {
	if !n.IsDir() {
		return newFileInode(n.store, n.Path, n.Value, n.Parent)
	}

	d := newDirInode(n.store, n.Path, n.Parent)
	for key, child := range n.Children {
		d.Children[key] = child.Clone()
	}

	return d
}
