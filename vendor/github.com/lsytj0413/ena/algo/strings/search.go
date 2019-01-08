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

package strings

// Find will return the index where sub in str
// -1 will return when sub not in str
func Find(str string, sub string) int {
	return _kmpFind(str, sub)
}

func _kmpFind(str string, sub string) int {
	if sub == "" {
		return 0
	}

	_next := func(sub string) []int {
		if len(sub) == 0 {
			return []int{}
		}

		r := make([]int, len(sub))
		r[0] = -1
		k, i := -1, 0

		for i < len(sub)-1 {
			if k == -1 || sub[i] == sub[k] {
				i++
				k++
				if sub[i] == sub[k] {
					r[i] = r[k]
				} else {
					r[i] = k
				}
			} else {
				k = r[k]
			}
		}
		return r
	}
	next := _next(sub)
	i, j := 0, 0
	for i < len(str) && j < len(sub) {
		if j == -1 || str[i] == sub[j] {
			i++
			j++
		} else {
			j = next[j]
		}
	}

	if j == len(sub) {
		return i - j
	}

	return -1
}

func _normalFind(str string, sub string) int {
	for i := 0; i < len(str)-len(sub)+1; i++ {
		j := 0
		for j = 0; j < len(sub); j++ {
			if str[i+j] != sub[j] {
				break
			}
		}
		if j == len(sub) {
			return i
		}
	}

	return -1
}
