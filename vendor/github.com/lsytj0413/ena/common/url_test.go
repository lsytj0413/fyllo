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

package common

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type urlTestSuite struct {
	suite.Suite
}

func (p *urlTestSuite) TestCleanURLPath() {
	var cleanTests = []struct {
		path, result string
	}{
		// Already clean
		{"/", "/"},
		{"/abc", "/abc"},
		{"/a/b/c", "/a/b/c"},
		{"/abc/", "/abc/"},
		{"/a/b/c/", "/a/b/c/"},

		// missing root
		{"", "/"},
		{"a/", "/a/"},
		{"abc", "/abc"},
		{"abc/def", "/abc/def"},
		{"a/b/c", "/a/b/c"},

		// Remove doubled slash
		{"//", "/"},
		{"/abc//", "/abc/"},
		{"/abc/def//", "/abc/def/"},
		{"/a/b/c//", "/a/b/c/"},
		{"/abc//def//ghi", "/abc/def/ghi"},
		{"//abc", "/abc"},
		{"///abc", "/abc"},
		{"//abc//", "/abc/"},

		// Remove . elements
		{".", "/"},
		{"./", "/"},
		{"/abc/./def", "/abc/def"},
		{"/./abc/def", "/abc/def"},
		{"/abc/.", "/abc/"},

		// Remove .. elements
		{"..", "/"},
		{"../", "/"},
		{"../../", "/"},
		{"../..", "/"},
		{"../../abc", "/abc"},
		{"/abc/def/ghi/../jkl", "/abc/def/jkl"},
		{"/abc/def/../ghi/../jkl", "/abc/jkl"},
		{"/abc/def/..", "/abc"},
		{"/abc/def/../..", "/"},
		{"/abc/def/../../..", "/"},
		{"/abc/def/../../..", "/"},
		{"/abc/def/../../../ghi/jkl/../../../mno", "/mno"},

		// Combinations
		{"abc/./../def", "/def"},
		{"abc//./../def", "/def"},
		{"abc/../../././../def", "/def"},
	}

	for _, v := range cleanTests {
		p.Equal(v.result, CleanURLPath(v.path))
	}
}

func TestURLTestSuite(t *testing.T) {
	p := &urlTestSuite{}
	suite.Run(t, p)
}
