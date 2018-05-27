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

package conf

import "context"

// Config is iploc server config instance
type Config struct {
	Name                   string `json:"name"`
	DefaultListenClientURL string `json:"listenClientUrl"`
	IsDebug                bool
	IsPprof                bool

	// snowflake config
	SnowflakeConfig GeneratorConfig
	Snowflake       SnowflakeGenerator

	// segment config
	SegmentConfig GeneratorConfig
	Segment       SegmentGenerator

	// random config
	RandomConfig GeneratorConfig
	Random       RandomGenerator

	// 客户端证书
	ClientTLSInfo TLSInfo
	IsTLSEnable   bool
}

// GeneratorConfig is generator config
type GeneratorConfig struct {
	Plugin string
	Args   string
}

// RandomResult for RandomGenerator Next
type RandomResult struct {
	Identify string `json:"identify"`
	Next     string `json:"next"`
}

// RandomGenerator is defines for random
type RandomGenerator interface {
	Next(context.Context) (*RandomResult, error)
}

// SnowflakeResult for SnowflakeGenerator Next
type SnowflakeResult struct {
	Next      uint64 `json:"next"`
	Timestamp uint64 `json:"timestamp"`
	Sequence  uint64 `json:"sequence"`
	Machine   uint64 `json:"machine"`
	Tag       uint64 `json:"tag"`
}

// SnowflakeGenerator is defines for snowflake
type SnowflakeGenerator interface {
	Next(context.Context, uint64) (*SnowflakeResult, error)
}

// SegmentResult for SegmentGenerator Next
type SegmentResult struct {
	Next uint64 `json:"next"`
	Tag  string `json:"tag"`
}

// SegmentGenerator is defines for segment
type SegmentGenerator interface {
	Next(context.Context, string) (*SegmentResult, error)
}

// TLSInfo is tls certificate info
type TLSInfo struct {
	CertFile       string
	KeyFile        string
	TrustedCAFile  string
	ClientCertAuth bool
	CRLFile        string

	InsecureSkipVerify bool
}

const (
	defaultName = "fyllo"
)

// New will construct a Config instance
func New() *Config {
	c := &Config{
		Name: defaultName,
	}
	return c
}
