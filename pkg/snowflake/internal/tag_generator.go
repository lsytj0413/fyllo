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

	ierror "github.com/lsytj0413/fyllo/pkg/error"
)

// Generator implement bind tag
type tagGenerator struct {
	tag           uint64
	machine       uint64
	lastTimestamp uint64
	mutex         sync.Mutex

	s sequence
	t ms
}

func (g *tagGenerator) processWithTimestamp(now uint64) (uint64, error) {
	if now > g.lastTimestamp {
		g.lastTimestamp = now
		g.s.reset()
	} else if now == g.lastTimestamp {
		if g.s.isOutRange() {
			g.lastTimestamp = g.t.waitForNextMS(now)
			g.s.reset()
		}
	} else {
		return 0, ierror.NewError(ierror.EcodeTimestampRewind, fmt.Sprintf("current timestamp[%d] less than last timestamp[%d]", now, g.lastTimestamp))
	}

	return g.lastTimestamp, nil
}

func (g *tagGenerator) Next(tag uint64) (*Result, error) {
	if tag != g.tag {
		panic("generator tag doesn't equal to args tag")
	}

	g.mutex.Lock()
	defer g.mutex.Unlock()

	_, err := g.processWithTimestamp(g.t.currMicroSeconds())
	if err != nil {
		return nil, err
	}

	r := &Result{
		Timestamp: g.lastTimestamp,
		Tag:       g.tag,
		Machine:   g.machine,
	}
	r.Sequence, err = g.s.next()
	if err != nil {
		return nil, err
	}
	r.Next = newID(r.Timestamp, r.Machine, r.Tag, r.Sequence)

	return r, nil
}

func newID(t uint64, m uint64, b uint64, s uint64) uint64 {
	timestamp := ((t - StandardTimestamp) & TimestampMask) << TimestampShiftNum
	maskMachineID := (m & MachineIDMask) << MachineIDShiftNum
	bid := (t & BusIDMask) << BusIDShiftNum
	sid := (s & SerialIDMask)

	return (timestamp | maskMachineID | bid | sid)
}
