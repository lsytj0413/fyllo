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

// CleanURLPath is the URL version of path.Clean, it returns a canonical URL path
// for p, eliminating . and .. elements.
//
// The following rules are applied iteratively until no further processing can
// be done:
//	1. Replace multiple slashes with a single slash.
//	2. Eliminate each . path name element (the current directory).
//	3. Eliminate each inner .. path name element (the parent directory)
//	   along with the non-.. element that precedes it.
//	4. Eliminate .. elements that begin a rooted path:
//	   that is, replace "/.." by "/" at the beginning of a path.
//
// If the result of this process is an empty string, "/" is returned
// copy from: httprouter.CleanPath
func CleanURLPath(p string) string {
	if p == "" {
		return "/"
	}

	n := len(p)
	// pre allocate, +1 for p not start with /
	var buf = make([]byte, n+1)

	// r is reading index from path, next byte to process
	r := 1
	// w is writing to buf, next byte to write, why start at 1? because buf[0] always is '/'
	w := 1

	// promise return path start with /
	// need specical process because at the loop will write / before every segment start,
	if p[0] != '/' {
		// reset r to 0, next byte is index 0, else r is 1 and jump over it
		r = 0
	}
	buf[0] = '/'

	// trailing is indicate should append / at last, > 1 for only '/'
	trailing := n > 1 && p[n-1] == '/'

	// loop for p
	for r < n {
		switch {
		case p[r] == '/':
			// jump over it, slash is added at segment start
			r++
		case p[r] == '.' && r+1 == n:
			// last char is ., should added slash at last
			trailing = true
			r++
		case p[r] == '.' && p[r+1] == '/':
			// ./ element, jump over
			r += 2
		case p[r] == '.' && p[r+1] == '.' && (r+2 == n || p[r+2] == '/'):
			// .. or ../
			r += 3

			if w > 1 { // jump over last '/'
				w--
				for w > 1 && buf[w] != '/' {
					w--
				}
			}
		default:
			if w > 1 { // if has write after first '/'
				buf[w] = '/'
				w++
			}

			for r < n && p[r] != '/' {
				buf[w] = p[r]
				w++
				r++
			}
		}
	}

	if trailing && w > 1 {
		buf[w] = '/'
		w++
	}

	rbuf := make([]byte, w)
	copy(rbuf, buf)
	return string(rbuf)
}
