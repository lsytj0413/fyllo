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
	"testing"

	"github.com/lsytj0413/ena/cerror"
	"github.com/stretchr/testify/suite"
)

type inodeTestSuite struct {
	suite.Suite
}

func (s *inodeTestSuite) TestIsHidden() {
	values := map[string]struct {
		node     *inode
		isHidden bool
	}{
		"1": {node: newFileInode(nil, "/", "", nil), isHidden: false},
		"2": {node: newFileInode(nil, "/xx", "", nil), isHidden: false},
		"3": {node: newFileInode(nil, "/.key", "", nil), isHidden: true},
		"4": {node: newFileInode(nil, "/key/xx", "", nil), isHidden: false},
		"5": {node: newFileInode(nil, "/key/xx/.ky", "", nil), isHidden: true},
		"6": {node: newDirInode(nil, "/", nil), isHidden: false},
		"7": {node: newDirInode(nil, "/.k", nil), isHidden: true},
		"8": {node: newDirInode(nil, "/.k/v", nil), isHidden: false},
		"9": {node: newDirInode(nil, "/.k/xx/.ee", nil), isHidden: true},
	}

	for _, v := range values {
		s.Equal(v.isHidden, v.node.IsHidden())
	}
}

func (s *inodeTestSuite) TestIsDir() {
	values := map[string]struct {
		node     *inode
		isHidden bool
	}{
		"1": {node: newFileInode(nil, "/", "", nil), isHidden: false},
		"2": {node: newFileInode(nil, "/xx", "", nil), isHidden: false},
		"3": {node: newFileInode(nil, "/.key", "", nil), isHidden: false},
		"4": {node: newFileInode(nil, "/key/xx", "", nil), isHidden: false},
		"5": {node: newFileInode(nil, "/key/xx/.ky", "", nil), isHidden: false},
		"6": {node: newDirInode(nil, "/", nil), isHidden: true},
		"7": {node: newDirInode(nil, "/.k", nil), isHidden: true},
		"8": {node: newDirInode(nil, "/.k/v", nil), isHidden: true},
		"9": {node: newDirInode(nil, "/.k/xx/.ee", nil), isHidden: true},
	}

	for _, v := range values {
		s.Equal(v.isHidden, v.node.IsDir())
	}
}

func (s *inodeTestSuite) TestReadOk() {
	defaultValue := "test"
	node := newFileInode(nil, "/test", defaultValue, nil)
	v, err := node.Read()
	s.NoError(err)
	s.Equal(defaultValue, v)
}

func (s *inodeTestSuite) TestReadDirError() {
	node := newDirInode(nil, "/test", nil)
	v, err := node.Read()
	s.Equal("", v)
	s.Error(err)

	e := err.(*cerror.Error)
	s.Equal(EcodeNotFile, e.ErrorCode)
}

func (s *inodeTestSuite) TestWriteOk() {
	defaultValue := "test"
	node := newFileInode(nil, "/test", "", nil)
	err := node.Write(defaultValue)
	s.NoError(err)

	s.Equal(defaultValue, node.Value)
}

func (s *inodeTestSuite) TestWriteDirError() {
	node := newDirInode(nil, "/test", nil)
	err := node.Write("")
	s.Error(err)

	e := err.(*cerror.Error)
	s.Equal(EcodeNotFile, e.ErrorCode)
}

func (s *inodeTestSuite) TestAddOk() {
	dnode := newDirInode(nil, "/test", nil)
	fnode := newFileInode(nil, "/test/tt", "", dnode)

	err := dnode.Add(fnode)
	s.NoError(err)

	s.Equal(1, len(dnode.Children))
	child, err := dnode.GetChild("tt")
	s.NoError(err)
	s.Equal(fnode, child)

	fnode = newFileInode(nil, "/test/tt1", "", dnode)
	err = dnode.Add(fnode)
	s.NoError(err)

	s.Equal(2, len(dnode.Children))
	child, err = dnode.GetChild("tt1")
	s.NoError(err)
	s.Equal(fnode, child)
}

func (s *inodeTestSuite) TestAddNotDirError() {
	node := newFileInode(nil, "/test", "", nil)
	err := node.Add(node)

	s.Error(err)

	e := err.(*cerror.Error)
	s.Equal(EcodeNotDir, e.ErrorCode)
}

func (s *inodeTestSuite) TestAddExistsError() {
	dnode := newDirInode(nil, "/test", nil)
	fnode := newFileInode(nil, "/test/tt", "", dnode)

	err := dnode.Add(fnode)
	s.NoError(err)

	err = dnode.Add(fnode)
	s.Error(err)

	e := err.(*cerror.Error)
	s.Equal(EcodeExists, e.ErrorCode)
}

func (s *inodeTestSuite) TestGetChildOk() {
	dnode := newDirInode(nil, "/test", nil)
	fnode := newFileInode(nil, "/test/tt", "", dnode)

	err := dnode.Add(fnode)
	s.NoError(err)

	child, err := dnode.GetChild("tt")
	s.NoError(err)
	s.Equal(fnode, child)

	fnode = newFileInode(nil, "/test/tt1", "", dnode)
	err = dnode.Add(fnode)
	s.NoError(err)

	child, err = dnode.GetChild("tt1")
	s.NoError(err)
	s.Equal(fnode, child)
}

func (s *inodeTestSuite) TestGetChildNotDirError() {
	node := newFileInode(nil, "/test", "", nil)
	child, err := node.GetChild("test")

	s.Error(err)
	s.Nil(child)

	e := err.(*cerror.Error)
	s.Equal(EcodeNotDir, e.ErrorCode)
}

func (s *inodeTestSuite) TestGetChildNotExistsError() {
	dnode := newDirInode(nil, "/test", nil)
	fnode := newFileInode(nil, "/test/tt", "", dnode)

	err := dnode.Add(fnode)
	s.NoError(err)

	child, err := dnode.GetChild("tt1")
	s.Error(err)
	s.Nil(child)

	e := err.(*cerror.Error)
	s.Equal(EcodeNotExists, e.ErrorCode)
}

func (s *inodeTestSuite) TestListOk() {
	dnode := newDirInode(nil, "/test", nil)

	r, err := dnode.List()
	s.NoError(err)
	s.NotNil(r)
	s.Equal(0, len(r))

	fnode := newFileInode(nil, "/test/tt", "", dnode)
	err = dnode.Add(fnode)
	s.NoError(err)

	r, err = dnode.List()
	s.NoError(err)
	s.NotNil(r)
	s.Equal(1, len(r))
	s.Equal(fnode, r[0])
}

func (s *inodeTestSuite) TestListNotDirError() {
	fnode := newFileInode(nil, "/test/tt", "", nil)
	r, err := fnode.List()
	s.Nil(r)
	s.Error(err)

	e := err.(*cerror.Error)
	s.Equal(EcodeNotDir, e.ErrorCode)
}

func (s *inodeTestSuite) TestRemoveFileOk() {
	dnode := newDirInode(nil, "/test", nil)
	fnode := newFileInode(nil, "/test/tt", "", dnode)

	err := dnode.Add(fnode)
	s.NoError(err)

	err = fnode.Remove(false, false)
	s.NoError(err)

	s.Equal(0, len(dnode.Children))
}

func (s *inodeTestSuite) TestRemoveDirOk() {
	dnode := newDirInode(nil, "/test", nil)
	dnode2 := newDirInode(nil, "/test/tt", dnode)

	err := dnode.Add(dnode2)
	s.NoError(err)

	err = dnode2.Remove(true, false)
	s.NoError(err)

	s.Equal(0, len(dnode.Children))
}

func (s *inodeTestSuite) TestRemoveDirRecursiveOk() {
	dnode := newDirInode(nil, "/test", nil)
	dnode2 := newDirInode(nil, "/test/tt", dnode)

	err := dnode.Add(dnode2)
	s.NoError(err)

	fnode := newFileInode(nil, "/test/tt/tt", "", dnode2)

	err = dnode2.Add(fnode)
	s.NoError(err)

	err = dnode2.Remove(true, true)
	s.NoError(err)

	s.Equal(0, len(dnode.Children))
	s.Equal(0, len(dnode2.Children))
}

func (s *inodeTestSuite) TestRemoveNotFileError() {
	dnode := newDirInode(nil, "/test", nil)
	dnode2 := newDirInode(nil, "/test/tt", dnode)

	err := dnode.Add(dnode2)
	s.NoError(err)

	err = dnode2.Remove(false, false)
	s.Error(err)

	e := err.(*cerror.Error)
	s.Equal(EcodeNotFile, e.ErrorCode)
}

func (s *inodeTestSuite) TestRemoveDirNotEmptyError() {
	dnode := newDirInode(nil, "/test", nil)
	dnode2 := newDirInode(nil, "/test/tt", dnode)

	err := dnode.Add(dnode2)
	s.NoError(err)

	fnode := newFileInode(nil, "/test/tt/tt", "", dnode2)

	err = dnode2.Add(fnode)
	s.NoError(err)

	err = dnode2.Remove(true, false)
	s.Error(err)

	e := err.(*cerror.Error)
	s.Equal(EcodeDirNotEmpty, e.ErrorCode)
}

func (s *inodeTestSuite) TestCloneFileOk() {
	fnode := newFileInode(nil, "/test/tt/tt", "", nil)
	n := fnode.Clone()

	s.Equal(*fnode, *n)
}

func (s *inodeTestSuite) TestCloneDirOk() {
	dnode := newDirInode(nil, "/test", nil)
	dnode2 := newDirInode(nil, "/test/tt", dnode)

	err := dnode.Add(dnode2)
	s.NoError(err)

	fnode := newFileInode(nil, "/test/tt/tt", "", dnode2)

	err = dnode2.Add(fnode)
	s.NoError(err)

	n := dnode.Clone()
	s.Equal(*dnode, *n)
}

func TestInodeTestSuite(t *testing.T) {
	s := &inodeTestSuite{}
	suite.Run(t, s)
}
