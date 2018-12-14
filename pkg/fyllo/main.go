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
	"net/url"

	"github.com/lsytj0413/ena/logger"

	randomBuilder "github.com/lsytj0413/fyllo/pkg/random/builder"
	segmentBuilder "github.com/lsytj0413/fyllo/pkg/segment/builder"
	"github.com/lsytj0413/fyllo/pkg/server"
	snowflakeBuilder "github.com/lsytj0413/fyllo/pkg/snowflake/builder"
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

	randomProviderName = flag.String("random-provider", "uuid", fmt.Sprintf("Random provider type. Available values: %s", randomBuilder.AvailableProvidersDescription))
	randomProviderArgs = flag.String("random-provider-args", "", "Random provider external arguments.")

	segmentProviderName = flag.String("segment-provider", "mem", fmt.Sprintf("Segment provider type. Available values: %s", segmentBuilder.AvailableProvidersDescription))
	segmentProviderArgs = flag.String("segment-provider-args", "", "Segment provider external arguments.")

	snowflakeProviderName = flag.String("snowflake-provider", "rock", fmt.Sprintf("Snowflake provider type. Available values: %s", snowflakeBuilder.AvailableProvidersDescription))
	snowflakeProviderArgs = flag.String("snowflake-provider-args", "0", "Snowflake provider external arguments.")
)

// Main is entrance for follymain application
func Main() {
	flag.Parse()

	if 0 != len(flag.Args()) {
		logger.Errorf("'%s' is not a valid flag", flag.Arg(0))
		return
	}

	if *isDebug {
		logger.SetLogLevel(logger.DebugLevel)
	}

	snowflakeProvider, err := buildSnowflakeProvider(*snowflakeProviderName, *snowflakeProviderArgs)
	if err != nil {
		logger.Errorf("build snowflake provider[name=%s][args=%s] failed, %v", *snowflakeProviderName, *snowflakeProviderArgs, err)
		return
	}

	segmentProvider, err := buildSegmentProvider(*segmentProviderName, *segmentProviderArgs)
	if err != nil {
		logger.Errorf("build segment provider[name=%s][args=%s] failed, %v", *segmentProviderName, *segmentProviderArgs, err)
		return
	}

	randomProvider, err := buildRandomProvider(*randomProviderName, *randomProviderArgs)
	if err != nil {
		logger.Errorf("build random provider[name=%s][args=%s] failed, %v", *randomProviderName, *randomProviderArgs, err)
		return
	}

	clientURL, err := url.Parse(*listenClientURL)
	if err != nil {
		logger.Errorf("%s is not valid url: %v", *listenClientURL, err)
		return
	}

	option := &server.Options{
		Name:                   *name,
		ListenClientURL:        clientURL,
		ClientCertFile:         *clientCertFile,
		ClientKeyFile:          *clientKeyFile,
		ClientTrustedCAFile:    *clientTrustedCAFile,
		IsClientCertAuthEnable: *isClientCertAuthEnable,
		IsInsecureSkipVerify:   *isInsecureSkipVerify,
		ClientCrlFile:          *clientCrlFile,
		IsPprof:                *isPprof,
	}
	srv, err := server.NewServer(option)
	if err != nil {
		logger.Errorf("NewServer failed: %v", err)
		return
	}

	installer := server.NewInstallers(
		&versionService{},
		&snowflakeService{
			provider: snowflakeProvider,
		},
		&segmentService{
			provider: segmentProvider,
		},
		&randomService{
			provider: randomProvider,
		},
	)

	stop := make(chan struct{})
	err = srv.Start(installer, stop)
	if err != nil {
		logger.Errorf("Server Start failed: %v", err)
		return
	}

	<-stop
	logger.Infof("Server stop, exit")
}
