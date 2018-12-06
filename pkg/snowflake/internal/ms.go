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
	"time"
)

const (
	// Nano2MicroRatio is the ratio convert nanoseconds to microseconds
	Nano2MicroRatio = 1000000
)

var (
	nano func() uint64
)

type ms interface {
	currMicroSeconds() uint64
	waitForNextMS(uint64) uint64
}

type defMs struct {
}

func (m *defMs) currMicroSeconds() uint64 {
	return nano() / Nano2MicroRatio
}

func (m *defMs) waitForNextMS(timestamp uint64) uint64 {
	for {
		now := m.currMicroSeconds()
		if now > timestamp {
			return now
		}
	}
}

func newMs() ms {
	return &defMs{}
}

func init() {
	nano = func() uint64 {
		return uint64(time.Now().UnixNano())
	}
}