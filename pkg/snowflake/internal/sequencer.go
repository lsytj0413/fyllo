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
	"sync"
	"sync/atomic"

	"github.com/lsytj0413/fyllo/pkg/snowflake"

	ierror "github.com/lsytj0413/fyllo/pkg/error"
)

// Sequencer for produce snowflake sequence number
type Sequencer interface {
	Next() (sequenceNumber uint64, timestamp uint64, err error)
}

type sequencer struct {
	now      uint64
	lastTick uint64
	mutex    sync.Mutex

	t Ticker
}

func (s *sequencer) Next() (sequenceNumber uint64, timestamp uint64, err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	err = s.checkTick()
	if err != nil {
		return
	}

	sequenceNumber = atomic.AddUint64(&s.now, 1)
	timestamp = s.lastTick
	return
}

func (s *sequencer) checkTick() error {
	n := s.t.Current()

	switch {
	case n > s.lastTick:
		s.lastTick = n
		s.now = 0
	case n == s.lastTick:
		if s.now >= snowflake.MaxSequenceValue {
			s.lastTick = s.t.WaitForNext(n)
			s.now = 0
		}
	default:
		return ierror.NewError(ierror.EcodeTimestampRewind, fmt.Sprintf("current timestamp[%d] less than last timestamp[%d]", n, s.lastTick))
	}
	return nil
}

// NewSequencer return Sequencer instance
func NewSequencer() Sequencer {
	return &sequencer{
		now:      0,
		lastTick: 0,
		t:        defaultTicker,
	}
}
