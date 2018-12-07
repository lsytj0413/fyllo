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

import "time"

const (
	// Nano2MicroRatio is the ratio convert nanoseconds to microseconds
	Nano2MicroRatio = 1000000
)

var (
	nano func() uint64
)

// Ticker for snowflake tick
type Ticker interface {
	Current() uint64
	WaitForNext(uint64) uint64
}

type ticker struct {
}

func (t *ticker) Current() uint64 {
	return nano() / Nano2MicroRatio
}

func (t *ticker) WaitForNext(cur uint64) uint64 {
	for {
		now := t.Current()
		if now > cur {
			return now
		}
	}
}

var (
	defaultTicker Ticker = &ticker{}
)

func init() {
	nano = func() uint64 {
		return uint64(time.Now().UnixNano())
	}
}
