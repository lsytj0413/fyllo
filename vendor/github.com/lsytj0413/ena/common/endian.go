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

package common

import (
	"encoding/binary"
	"unsafe"
)

const intWidth int = int(unsafe.Sizeof(0))

var byteOrder binary.ByteOrder

// ByteOrder returns the byte order for the CPU's native endianness.
func ByteOrder() binary.ByteOrder {
	return byteOrder
}

func init() {
	var i = 0x01
	if v := (*[intWidth]byte)(unsafe.Pointer(&i)); v[0] == 0 {
		byteOrder = binary.BigEndian
	} else {
		byteOrder = binary.LittleEndian
	}
}
