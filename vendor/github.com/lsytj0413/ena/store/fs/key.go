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
	"path"
	"strings"
)

const (
	root = "/"
)

// construct a nodePath to a key
func key(nodePath string) string {
	return path.Clean(path.Join(root, nodePath))
}

// construct a nodepath from dirname and filename
func keyFromDirAndFile(dir string, file string) string {
	return path.Clean(path.Join(key(dir), file))
}

// check nodePath is root key
func isRoot(nodePath string) bool {
	return nodePath == root
}

// split nodePath to array of components
func components(nodePath string) []string {
	return strings.Split(key(nodePath), root)
}

// name returns node filename component
func name(nodePath string) string {
	_, name := split(nodePath)
	return name
}

// split nodePath to dirname and filename
func split(nodePath string) (string, string) {
	return path.Split(key(nodePath))
}
