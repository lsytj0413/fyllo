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

func maxUint(x uint, y uint) uint {
	if x < y {
		return y
	}

	return x
}

// IsJumpEnd returns is able jump to end
func IsJumpEnd(steps []uint) bool {
	if len(steps) == 0 || len(steps) == 1 {
		return true
	}

	i, max := 0, uint(0)

	for ; uint(i) <= max && i < len(steps); i++ {
		max = maxUint(max, steps[i]+uint(i))
	}

	return i >= len(steps)
}

// JumpShortedSteps returns the shorted steps jump to end
func JumpShortedSteps(steps []uint) int {
	if len(steps) == 0 {
		return -1
	}

	if len(steps) == 1 {
		return 0
	}

	return jumpShortedStepsImpl(steps, 0)
}

func jumpShortedStepsImpl(steps []uint, i int) int {
	if i >= len(steps) {
		return 1
	}
	if i == len(steps)-1 {
		return 0
	}

	nextJump := i + int(steps[i])
	if nextJump >= len(steps) {
		return 1
	}

	maxJump := nextJump
	for j := i + 1; j < nextJump; j++ {
		if int(steps[j])+j > maxJump {
			maxJump = j
		}
	}

	return 1 + jumpShortedStepsImpl(steps, maxJump)
}
