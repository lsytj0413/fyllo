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

import "fmt"

// @struct serial
// @brief 序列号产生器
type serial uint64

const (
	// MaxSerialNumber @var
	// @brief 最大可用序列号
	MaxSerialNumber serial = (1 << 10) - 1
)

// get the next serial at this microsecond
func (p *serial) next() (uint64, error) {
	if *p >= MaxSerialNumber {
		panic(fmt.Errorf("out of serial"))
	}

	sNumber := *p
	*p++

	return uint64(sNumber), nil
}

// reset serial, at every microsecond should reset serail first
func (p *serial) reset() {
	*p = 0
}
