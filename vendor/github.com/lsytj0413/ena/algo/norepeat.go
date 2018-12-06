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

// FindNoRepeat will find the longest norepeat substring
func FindNoRepeat(src string) (string, int) {
	if src == "" {
		return src, 0
	}

	return findNoRepeatImpl([]rune(src))
}

func findNoRepeatImpl(src []rune) (string, int) {
	startIndex := 0
	maxLength := 1
	maxStartIndex := 0
	cachedIndex := make(map[rune]int, len(src))
	cachedIndex[src[0]] = 0

	for i := 1; i < len(src); i++ {

		lastIndex, exists := cachedIndex[src[i]]
		if exists && lastIndex >= startIndex {
			// 出现在当前连续最长中, 则更新 startIndex
			startIndex = lastIndex + 1
		}

		// 判断当前的是否大于历史
		currLentgh := i - startIndex + 1
		if currLentgh > maxLength {
			maxStartIndex = startIndex
			maxLength = currLentgh
		}
		cachedIndex[src[i]] = i
	}

	return string(src[maxStartIndex : maxStartIndex+maxLength]), maxLength
}
