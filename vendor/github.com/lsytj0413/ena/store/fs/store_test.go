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

	"github.com/stretchr/testify/suite"
)

type defFileSystemStoreTestSuite struct {
	suite.Suite
	store *defFileSystemStore
}

func (s *defFileSystemStoreTestSuite) SetupTest() {
	s.store = newDefFileSystemStore()
}

func (s *defFileSystemStoreTestSuite) TearDownTest() {
	s.store = nil
}

func (s *defFileSystemStoreTestSuite) TestSetOk() {
	r, err := s.store.Set("xxx", false, "xxx")
	s.Suite.NoError(err)

	s.Equal(Set, r.Action)
	s.Equal(false, r.CurrNode.Dir)
	s.Equal("/xxx", r.CurrNode.Key)
	s.Equal("xxx", *(r.CurrNode.Value))
	s.Nil(r.PrevNode)
}

func (s *defFileSystemStoreTestSuite) TestGetOk() {
	_, err := s.store.Set("xxx", false, "xxx")
	s.Suite.NoError(err)

	r, err := s.store.Get("xxx", false, false)
	s.NoError(err)
	s.Equal(Get, r.Action)
	s.Equal(false, r.CurrNode.Dir)
	s.Equal("/xxx", r.CurrNode.Key)
	s.Equal("xxx", *(r.CurrNode.Value))
	s.Nil(r.PrevNode)
}

func (s *defFileSystemStoreTestSuite) TestUpdateOk() {
	_, err := s.store.Set("xxx", false, "xxx")
	s.Suite.NoError(err)

	r, err := s.store.Update("/xxx", "newxxx")
	s.NoError(err)
	s.Equal(Update, r.Action)
	s.Equal(false, r.CurrNode.Dir)
	s.Equal("/xxx", r.CurrNode.Key)
	s.Equal("newxxx", *(r.CurrNode.Value))

	s.Equal(false, r.PrevNode.Dir)
	s.Equal("/xxx", r.PrevNode.Key)
	s.Equal("xxx", *(r.PrevNode.Value))
}

func (s *defFileSystemStoreTestSuite) TestCreateOk() {
	r, err := s.store.Create("xxx", false, "xxx")
	s.Suite.NoError(err)

	s.Equal(Create, r.Action)
	s.Equal(false, r.CurrNode.Dir)
	s.Equal("/xxx", r.CurrNode.Key)
	s.Equal("xxx", *(r.CurrNode.Value))
	s.Nil(r.PrevNode)
}

func (s *defFileSystemStoreTestSuite) TestDeleteOk() {
	_, err := s.store.Create("xxx", false, "xxx")
	s.NoError(err)

	r, err := s.store.Delete("xxx", false, false)
	s.NoError(err)

	s.Equal(Delete, r.Action)
	s.Equal(false, r.PrevNode.Dir)
	s.Equal("/xxx", r.PrevNode.Key)
	s.Equal("xxx", *(r.PrevNode.Value))
	s.Equal(false, r.CurrNode.Dir)
	s.Equal("/xxx", r.CurrNode.Key)
	s.Equal("xxx", *(r.CurrNode.Value))
}

func TestDefFilesystemStoreTestSuite(t *testing.T) {
	s := &defFileSystemStoreTestSuite{}
	suite.Run(t, s)
}
