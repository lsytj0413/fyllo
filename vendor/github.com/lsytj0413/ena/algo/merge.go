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

package algo

// Merge two sorted int slice
func Merge(a1 []int, a2 []int) []int {
	r := make([]int, len(a1)+len(a2))
	i, j := 0, 0

	for i < len(a1) && j < len(a2) {
		v1, v2 := a1[i], a2[j]
		if v1 < v2 {
			r[i+j] = v1
			i++
		} else {
			r[i+j] = v2
			j++
		}
	}

	for i < len(a1) {
		r[i+j] = a1[i]
		i++
	}
	for j < len(a2) {
		r[i+j] = a2[j]
		j++
	}

	return r
}
