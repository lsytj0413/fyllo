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

package internal

import (
	"fmt"

	ierror "github.com/lsytj0413/fyllo/pkg/error"
)

type sequence interface {
	next() (uint64, error)
	reset()
	isOutRange() bool
}

type defSequence struct {
	max uint64
	now uint64
}

func (s *defSequence) next() (uint64, error) {
	if s.now >= s.max {
		return 0, ierror.NewError(ierror.EcodeSequenceOutOfRange,
			fmt.Sprintf("sequence should in range of [0, %d)", s.max))
	}

	tmp := s.now
	s.now++
	return tmp, nil
}

func (s *defSequence) reset() {
	s.now = 0
}

func (s *defSequence) isOutRange() bool {
	return s.now >= s.max
}

func newSequence(max uint64) sequence {
	return &defSequence{
		max: max,
		now: 0,
	}
}
