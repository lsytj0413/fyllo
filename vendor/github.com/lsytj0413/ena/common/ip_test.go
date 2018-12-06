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
	"testing"

	"github.com/stretchr/testify/suite"
)

type ipTestSuite struct {
	suite.Suite
}

func (p *ipTestSuite) TestIp2Uint32Ok() {
	values := map[string]uint32{
		"0.0.0.0":         0,
		"0.255.255.255":   16777215,
		"1.0.0.255":       16777471,
		"1.0.127.255":     16809983,
		"1.1.128.0":       16875520,
		"1.8.1.255":       17302015,
		"1.11.0.0":        17498112,
		"36.16.214.255":   605083391,
		"36.43.28.0":      606804992,
		"58.35.247.255":   975435775,
		"61.155.235.67":   1033628483,
		"116.113.210.138": 1953616522,
		"185.60.104.0":    3107743744,
		"202.147.6.0":     3398632960,
		"223.255.255.0":   3758096128,
		"255.255.255.255": 4294967295,
	}
	for k, v := range values {
		actual, err := Ip2Uint32(k)
		p.Equal(err, nil)
		p.Equal(v, actual)
	}
}

func (p *ipTestSuite) TestIp2Uint32FailInvalid() {
	values := []string{"ip",
		"0.0.0.-1",
		"-1.0.0.0",
		"255.255.255.-1",
		"255.255.255.256",
		"256.255.255.255",
		"0.0.0.0.",
		"0.0.0.0.0",
		"0..0.0.0",
		".0.0.0",
		"0.0.0.",
		"CDCD:910A:2222:5498:8475:1111:3900:",
		"CDCD:910A:2222:5498:8475:1111:",
		":ffff:192.168.89.9",
		":192.168.89.9",
	}
	for _, value := range values {
		_, err := Ip2Uint32(value)
		p.EqualError(err, fmt.Sprintf("Invalid IP address: %s", value))
	}
}

func (p *ipTestSuite) TestIp2Uint32FailIpv6() {
	values := []string{"CDCD:910A:2222:5498:8475:1111:3900:2020",
		"1030::C9B4:FF12:48AA:1A2B",
		"2000:0:0:0:0:0:0:1",
		// blow ipv6 will be tranform to ipv4
		// "::ffff:192.168.89.9",
		// "::192.168.89.9",
	}
	for _, value := range values {
		_, err := Ip2Uint32(value)
		p.EqualError(err, fmt.Sprintf("Only Support IPv4 address: %s", value))
	}
}

func (p *ipTestSuite) TestUint322IpOK() {
	values := map[string]uint32{
		"0.0.0.0":         0,
		"0.255.255.255":   16777215,
		"1.0.0.255":       16777471,
		"1.0.127.255":     16809983,
		"1.1.128.0":       16875520,
		"1.8.1.255":       17302015,
		"1.11.0.0":        17498112,
		"36.16.214.255":   605083391,
		"36.43.28.0":      606804992,
		"58.35.247.255":   975435775,
		"61.155.235.67":   1033628483,
		"116.113.210.138": 1953616522,
		"185.60.104.0":    3107743744,
		"202.147.6.0":     3398632960,
		"223.255.255.0":   3758096128,
		"255.255.255.255": 4294967295,
	}
	for k, v := range values {
		actual := Uint322Ip(v)
		p.Equal(k, actual)
	}
}

func TestIpTestSuite(t *testing.T) {
	p := &ipTestSuite{}
	suite.Run(t, p)
}
