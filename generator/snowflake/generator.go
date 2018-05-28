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

package snowflake

import (
	"context"
	"fmt"
	"sync"

	"github.com/lsytj0413/fyllo/conf"
)

const (
	// MaxTag is maxinum of tag value
	MaxTag = 1 << 8

	// MaxMachine is maxinum of machine value
	MaxMachine = 1 << 4
)

const (
	// StandardTimestamp is the begining of timestamp
	StandardTimestamp uint64 = 1451606400000

	// TimestampMask is the mask of timestamp field
	TimestampMask uint64 = 0x000001FFFFFFFFFF

	// TimestampShiftNum is the shift number of timestamp field
	TimestampShiftNum = 22

	// MachineIDMask is the mask of machine field
	MachineIDMask uint64 = 0x000000000000000F

	// MachineIDShiftNum is the shift number of machine field
	MachineIDShiftNum = 18

	// BusIDMask is the mask of tag field
	BusIDMask uint64 = 0x00000000000000FF

	// BusIDShiftNum is the shift of tag field
	BusIDShiftNum = 10

	// SerialIDMask is the mask of sequence field
	SerialIDMask uint64 = 0x00000000000003FF
)

type generator struct {
	tag           uint64
	machine       uint64
	s             serial
	lastTimestamp uint64
	mutex         sync.Mutex
}

func (g *generator) next() (*conf.SnowflakeResult, error) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	_, err := g.processWithTimestamp(t.currMicroSeconds())
	if err != nil {
		return nil, err
	}

	r := &conf.SnowflakeResult{
		// Sequence:  g.s.next(),
		Timestamp: g.lastTimestamp,
		Machine:   g.machine,
		Tag:       g.tag,
	}
	r.Sequence, _ = g.s.next()
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

func (g *generator) processWithTimestamp(now uint64) (uint64, error) {
	if now > g.lastTimestamp {
		g.lastTimestamp = now
		g.s.reset()
	} else if now == g.lastTimestamp {
		if g.s >= MaxSerialNumber {
			g.lastTimestamp = t.waitForNextMS(now)
			g.s.reset()
		}
	} else {
		return 0, fmt.Errorf("current timestamp less than last timestamp")
	}

	return g.lastTimestamp, nil
}

type generatorMapper struct {
	mutex        sync.RWMutex
	machine      uint64
	tagGenerator map[uint64]*generator
}

var (
	mapper = generatorMapper{
		mutex:        sync.RWMutex{},
		machine:      0,
		tagGenerator: make(map[uint64]*generator, 32),
	}
)

func (m *generatorMapper) generator(tag uint64) *generator {
	getter := func() *generator {
		m.mutex.RLock()
		defer m.mutex.Unlock()

		g, ok := m.tagGenerator[tag]
		if !ok {
			return nil
		}

		return g
	}

	g := getter()
	if g != nil {
		return g
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	// try again
	g, ok := m.tagGenerator[tag]
	if ok {
		return g
	}

	g = &generator{
		tag: tag,
	}
	m.tagGenerator[tag] = g
	return g
}

func (m *generatorMapper) Next(c context.Context, tag uint64) (*conf.SnowflakeResult, error) {
	if err := checkTag(tag); err != nil {
		return nil, err
	}

	g := m.generator(tag)
	return g.next()
}

func checkTag(tag uint64) error {
	if tag >= MaxTag {
		return fmt.Errorf("invalid tag value=%d, it should be range of [0,%d)", tag, MaxTag)
	}

	return nil
}
