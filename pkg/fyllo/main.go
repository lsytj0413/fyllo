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

// Package fyllo is a distributed, unique id generation service.
package fyllo

import (
	"flag"
	"fmt"
)

var (
	name                   = flag.String("name", "", "Human-readable name for this member, if available")
	listenClientURL        = flag.String("listen-client-url", "http://localhost:80", "URL to listen on for client traffic.")
	clientCertFile         = flag.String("client-cert-file", "", "Path to the client server TLS cert file.")
	clientKeyFile          = flag.String("client-key-file", "", "Path to the client server TLS key file.")
	clientTrustedCAFile    = flag.String("client-trusted-ca-file", "", "Path th the client server TLS trusted CA cert file.")
	isClientCertAuthEnable = flag.Bool("client-cert-auth", false, "Enable client cert authentication.")
	isInsecureSkipVerify   = flag.Bool("client-auto-tls", false, "Client TLS using generated certificate.")
	clientCrlFile          = flag.String("client-crl-file", "", "Path to the client certificate revocation list file.")
	isDebug                = flag.Bool("debug", false, "Enable debug log output.")
	isPprof                = flag.Bool("pprof", false, "Enable pprof.")

	randomProviderName = flag.String("random-provider", "", "Random provider type. Available values: []")
	randomProviderArgs = flag.String("random-provider-args", "", "Random provider external arguments.")

	segmentProviderName = flag.String("segment-provider", "", "Segment provider type. Available values: []")
	segmentProviderArgs = flag.String("segment-provider-args", "", "Segment provider external arguments.")

	snowflakeProviderName = flag.String("snowflake-provider", "", "Snowflake provider type. Available values: []")
	snowflakeProviderArgs = flag.String("snowflake-provider-args", "", "Snowflake provider external arguments.")
)

// Main is entrance for follymain application
func Main() {
	fmt.Println("fyllo.Main")
}
