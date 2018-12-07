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

const (
	// MaxTag is maxinum of tag value
	MaxTag = 1 << 8

	// MaxMachine is maxinum of machine value
	MaxMachine = 1 << 4
)

const (
	// MaxSequenceValue is max value of sequence
	MaxSequenceValue uint64 = (1 << 10)
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

// MakeSnowflakeID generate snowflake ID from timestamp, machine id, bussiness id and sequence number.
// it use snowflake algorithm.
func MakeSnowflakeID(timestamp uint64, mid uint64, bid uint64, sequenceNumber uint64) uint64 {
	timestamp = ((timestamp - StandardTimestamp) & TimestampMask) << TimestampShiftNum
	mid = (mid & MachineIDMask) << MachineIDShiftNum
	bid = (bid & BusIDMask) << BusIDShiftNum
	sequenceNumber = (sequenceNumber & SerialIDMask)

	return (timestamp | mid | bid | sequenceNumber)
}
