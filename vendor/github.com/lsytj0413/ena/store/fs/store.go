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
	"errors"
	"fmt"
	"sync"
)

// Store defines a filesystem like kv store
type Store interface {
	// Get nodePath node infomation
	Get(nodePath string, recursive bool, sorted bool) (*Result, error)
	// Set value to nodePath
	Set(nodePath string, dir bool, value string) (*Result, error)
	// Update value to nodePath
	Update(nodePath string, newValue string) (*Result, error)
	// Create nodePath with value
	Create(nodePath string, dir bool, value string) (*Result, error)
	// Delete nodePath
	Delete(nodePath string, dir bool, recursive bool) (*Result, error)
	// Get the Stater Object
	Stater() Stater
	// Watch for key change event
	Watch(key string, recursive bool) (Watcher, error)
}

// defFileSystemStore implemented FileSystemStore interface
type defFileSystemStore struct {
	sync.RWMutex

	Root *inode

	stater     Stater
	watcherHub WatcherHub
}

// New creates a FileSystemStore with root directories
func New() Store {
	return newDefFileSystemStore()
}

func newDefFileSystemStore() *defFileSystemStore {
	s := new(defFileSystemStore)
	s.Root = newDirInode(s, "/", nil)
	s.stater = newStater()
	s.watcherHub = newWatchHub(1000)
	return s
}

// Stater returns the store.Stater object
func (s *defFileSystemStore) Stater() Stater {
	return s.stater
}

// Get returns Node which the nodePath specified
// If recursive is true, it will return all the content under the node path
func (s *defFileSystemStore) Get(nodePath string, recursive bool, sorted bool) (*Result, error) {
	s.RLock()
	defer s.RUnlock()

	var err error
	defer func() {
		if err == nil {
			s.stater.Inc(GetSuccess)
		} else {
			s.stater.Inc(GetFail)
		}
	}()

	n, err := s.get(nodePath)
	if err != nil {
		return nil, err
	}

	r := newResult(Get)
	r.CurrNode = inodeToNode(n, recursive, sorted)
	return r, nil
}

// Set create of replace the node at nodePath
func (s *defFileSystemStore) Set(nodePath string, dir bool, value string) (*Result, error) {
	s.Lock()
	defer s.Unlock()

	var err error
	defer func() {
		if err == nil {
			s.stater.Inc(SetSuccess)
		} else {
			s.stater.Inc(SetFail)
		}
	}()

	// First, get prevNode Value
	prevNode, err := s.get(nodePath)
	if err != nil && !Is(err, EcodeNotExists) {
		return nil, err
	}

	// remove exists inode before create
	if prevNode != nil {
		if prevNode.IsDir() {
			return nil, NewError(EcodeNotFile, fmt.Sprintf("Set %s: replace when prevNode is dir", nodePath))
		}

		err = prevNode.Remove(false, false)
		if err != nil {
			return nil, err
		}
	}

	n, err := s.create(nodePath, dir, value)
	if err != nil {
		return nil, err
	}

	e := newResult(Set)
	e.CurrNode = inodeToNode(n, false, false)
	if prevNode != nil {
		e.PrevNode = inodeToNode(prevNode, false, false)
	}

	// notify watcher
	s.watcherHub.notify(e)

	return e, nil
}

// Update updates the value of the node
// If the node is a directory, Update will fail
func (s *defFileSystemStore) Update(nodePath string, newValue string) (*Result, error) {
	s.Lock()
	defer s.Unlock()

	var err error
	defer func() {
		if err == nil {
			s.stater.Inc(UpdateSuccess)
		} else {
			s.stater.Inc(UpdateFail)
		}
	}()

	nodePath = key(nodePath)

	n, err := s.get(nodePath)
	if err != nil {
		return nil, err
	}

	if n.IsDir() {
		return nil, errors.New("Not file")
	}

	r := newResult(Update)
	r.PrevNode = inodeToNode(n, false, false)

	err = n.Write(newValue)
	if err != nil {
		return nil, err
	}

	r.CurrNode = inodeToNode(n, false, false)

	// notify watcher
	s.watcherHub.notify(r)

	return r, nil
}

// Create creates the node at nodePath.
// If the node has already exists, create will fail
// If any node on the path is file, create will fail
func (s *defFileSystemStore) Create(nodePath string, dir bool, value string) (*Result, error) {
	s.Lock()
	defer s.Unlock()

	var err error
	defer func() {
		if err == nil {
			s.stater.Inc(CreateSuccess)
		} else {
			s.stater.Inc(CreateFail)
		}
	}()

	n, err := s.create(nodePath, dir, value)
	if err != nil {
		return nil, err
	}

	e := newResult(Create)
	e.CurrNode = inodeToNode(n, false, false)

	s.watcherHub.notify(e)

	return e, nil
}

// create creates the node at nodePath
// If the node has already exists, fail with EcodeExists
// If any node on the path is file, fail with EcodeNotDir
func (s *defFileSystemStore) create(nodePath string, dir bool, value string) (*inode, error) {
	nodePath = key(nodePath)
	dirName, nodeName := split(nodePath)

	d, err := walk(dirName, s.Root, s.createDir)
	if err != nil {
		return nil, err
	}

	n, _ := d.GetChild(nodeName)
	if n != nil {
		return nil, NewError(EcodeExists, fmt.Sprintf("create %s", nodePath))
	}

	if !dir {
		n = newFileInode(s, nodePath, value, d)
	} else {
		n = newDirInode(s, nodePath, d)
	}
	d.Add(n)

	return n, nil
}

// createDir will auto create dirName under parent
// If is directory, return inode
// If does not exsits, create a new directory and return inode
// If is file and exists, return EcodeNotDir
func (s *defFileSystemStore) createDir(parent *inode, dirName string) (*inode, error) {
	node, ok := parent.Children[dirName]
	if ok {
		if node.IsDir() {
			return node, nil
		}

		return nil, NewError(EcodeNotDir, fmt.Sprintf("createDirResursive %s: parent=%s", dirName, parent.Path))
	}

	n := newDirInode(s, keyFromDirAndFile(parent.Path, dirName), parent)
	parent.Children[dirName] = n
	return n, nil
}

// Delete deletes the node at the given path
// If the node is a directory, recursive must be true to delete it
func (s *defFileSystemStore) Delete(nodePath string, dir bool, recursive bool) (*Result, error) {
	s.Lock()
	defer s.Unlock()

	var err error
	defer func() {
		if err == nil {
			s.stater.Inc(DeleteSuccess)
		} else {
			s.stater.Inc(DeleteFail)
		}
	}()

	nodePath = key(nodePath)

	if recursive {
		dir = true
	}

	n, err := s.get(nodePath)
	if err != nil {
		return nil, err
	}

	r := newResult(Delete)
	r.PrevNode = inodeToNode(n, false, false)
	r.CurrNode = inodeToNode(n, false, false)

	// BUG: notifyWatchers when delete children node
	err = n.Remove(dir, recursive)
	if err != nil {
		return nil, err
	}

	s.watcherHub.notify(r)

	return r, nil
}

func (s *defFileSystemStore) Watch(nodePath string, recursive bool) (Watcher, error) {
	s.RLock()
	defer s.RUnlock()

	nodePath = key(nodePath)
	w, err := s.watcherHub.watch(nodePath, recursive)
	if err != nil {
		return nil, err
	}

	return w, nil
}

// get find the nodePath inode
func (s *defFileSystemStore) get(nodePath string) (*inode, error) {
	nodePath = key(nodePath)

	walkFunc := func(parent *inode, name string) (*inode, error) {
		if !parent.IsDir() {
			return nil, NewError(EcodeNotDir, fmt.Sprintf("get %s: parent=%s", nodePath, parent.Path))
		}

		child, ok := parent.Children[name]
		if ok {
			return child, nil
		}

		return nil, NewError(EcodeNotExists, fmt.Sprintf("get %s", nodePath))
	}

	f, err := walk(nodePath, s.Root, walkFunc)
	if err != nil {
		return nil, err
	}

	return f, nil
}
