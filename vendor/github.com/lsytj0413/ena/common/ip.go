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
	"fmt"
	"net"
)

// Ip2Uint32 will encode ip string to uint32, only support ipv4
func Ip2Uint32(ip string) (uint32, error) {
	value := net.ParseIP(ip)
	if value == nil {
		return 0, fmt.Errorf("Invalid IP address: %s", ip)
	}

	// because ParseIP always return IPv6 IP, which len(value) is 16
	value = value.To4()
	if value == nil || 4 != len([]byte(value)) {
		return 0, fmt.Errorf("Only Support IPv4 address: %s", ip)
	}

	var dest uint32
	for i, seg := range value {
		shift := uint((3 - i) * 8)
		dest = dest | uint32(uint32(seg)<<shift)
	}

	return dest, nil
}

// Uint322Ip will encode ip uint32 to string, only support ipv4
func Uint322Ip(ip uint32) string {
	segs := IpSegments(ip)
	return fmt.Sprintf("%d.%d.%d.%d", segs[0], segs[1], segs[2], segs[3])
}

// IpSegments will split ip to byte array
func IpSegments(ip uint32) []uint8 {
	seg1 := uint8((ip & 0xFF000000) >> 24)
	seg2 := uint8((ip & 0x00FF0000) >> 16)
	seg3 := uint8((ip & 0x0000FF00) >> 8)
	seg4 := uint8((ip & 0x000000FF) >> 0)
	return []uint8{seg1, seg2, seg3, seg4}
}
