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

type keyTestSuite struct {
	suite.Suite
}

func (s *keyTestSuite) TestKeyOk() {
	values := map[string]string{
		"key":         "/key",
		"/key":        "/key",
		"key/key":     "/key/key",
		"/key/key":    "/key/key",
		"key//":       "/key",
		"//key//key":  "/key/key",
		"//key//key/": "/key/key",
		"key/":        "/key",
		`\key`:        "/\\key",
	}

	for k, v := range values {
		s.Equal(v, key(k))
	}
}

func (s *keyTestSuite) TestKeyFromDirAndFile() {
	values := map[string]struct {
		Dir  string
		File string
		Key  string
	}{
		"1": {Dir: "key", File: "key", Key: "/key/key"},
		"2": {Dir: "/key", File: "key", Key: "/key/key"},
		"3": {Dir: "/", File: "/key", Key: "/key"},
		"4": {Dir: "//key", File: "/key/", Key: "/key/key"},
	}

	for _, v := range values {
		s.Equal(v.Key, keyFromDirAndFile(v.Dir, v.File))
	}
}

func (s *keyTestSuite) TestIsRoot() {
	values := map[string]bool{
		"/":   true,
		".":   false,
		"//":  false,
		"xx":  false,
		"/.":  false,
		"/xx": false,
		"./":  false,
		"x/":  false,
	}

	for k, v := range values {
		s.Equal(v, isRoot(k))
	}
}

func (s *keyTestSuite) TestComponents() {
	values := map[string][]string{
		"/":             []string{"", ""},
		"/key":          []string{"", "key"},
		"/key/key/":     []string{"", "key", "key"},
		"//key/xxxx/xx": []string{"", "key", "xxxx", "xx"},
	}

	for k, v := range values {
		s.Equal(v, components(k))
	}
}

func (s *keyTestSuite) TestName() {
	values := map[string]string{
		"/":             "",
		"/key":          "key",
		"/key/key/":     "key",
		"//key/xxxx/xx": "xx",
	}

	for k, v := range values {
		s.Equal(v, name(k))
	}
}

func (s *keyTestSuite) TestSplit() {
	values := map[string]struct {
		Dir  string
		File string
		Key  string
	}{
		"1": {Dir: "/key/", File: "key", Key: "/key/key"},
		"2": {Dir: "/key/key/", File: "xxx", Key: "/key/key/xxx"},
		"3": {Dir: "/", File: "key", Key: "/key"},
		"4": {Dir: "/key/", File: "key", Key: "key/key"},
	}

	for _, v := range values {
		k1, v1 := split(v.Key)
		s.Equal(v.Dir, k1)
		s.Equal(v.File, v1)
	}
}

func TestKeyTestSuite(t *testing.T) {
	s := &keyTestSuite{}
	suite.Run(t, s)
}
