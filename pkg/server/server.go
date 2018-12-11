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

package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/lsytj0413/ena/logger"

	"github.com/lsytj0413/fyllo/pkg/random"
	randomBuilder "github.com/lsytj0413/fyllo/pkg/random/builder"
	"github.com/lsytj0413/fyllo/pkg/segment"
	segmentBuilder "github.com/lsytj0413/fyllo/pkg/segment/builder"
	"github.com/lsytj0413/fyllo/pkg/snowflake"
	snowflakeBuilder "github.com/lsytj0413/fyllo/pkg/snowflake/builder"
)

// Server is fyllo server interface
type Server interface {
	Start() (chan struct{}, error)
}

// Options for server option
type Options struct {
	Name                   string   `json:"name"`
	ListenClientURL        *url.URL `json:"url"`
	ClientCertFile         string   `json:"clientCertFile"`
	ClientKeyFile          string   `json:"clientKeyFile"`
	ClientTrustedCAFile    string   `json:"clientTrustedCAFile"`
	IsClientCertAuthEnable bool     `json:"isClientCertAuthEnable"`
	IsInsecureSkipVerify   bool     `json:"isInsecureSkipVerify"`
	ClientCrlFile          string   `json:"clientCrlFile"`
	IsDebug                bool     `json:"isDebug"`
	IsPprof                bool     `json:"isPprof"`
	IsTLSEnable            bool

	RandomProvider     string `json:"randomProvider"`
	RandomProviderArgs string `json:"randomProviderArgs"`

	SegmentProvider     string `json:"segmentProvider"`
	SegmentProviderArgs string `json:"segmentProviderArgs"`

	SnowflakeProvider     string `json:"snowflakeProvider"`
	SnowflakeProviderArgs string `json:"snowflakeProviderArgs"`
}

type server struct {
	option *Options

	randomProvider    random.Provider
	segmentProvider   segment.Provider
	snowflakeProvider snowflake.Provider

	stop chan struct{}
}

func (s *server) Start() (chan struct{}, error) {
	if s.option.IsDebug {
		logger.SetLogLevel(logger.DebugLevel)
	}

	if s.option.ListenClientURL.Scheme == "https" {
		if s.option.ClientCertFile == "" || s.option.ClientKeyFile == "" {
			return nil, fmt.Errorf("listen on https without keyfile or certfile exists")
		}

		s.option.IsTLSEnable = true
	}

	if s.option.IsClientCertAuthEnable {
		if s.option.ClientTrustedCAFile == "" && !s.option.IsInsecureSkipVerify {
			return nil, fmt.Errorf("client auth enable without client-trusted-ca-file or client-auto-tls set")
		}
	}

	tlsConfig := &tls.Config{}
	if s.option.IsClientCertAuthEnable {
		if s.option.IsInsecureSkipVerify {
			tlsConfig.InsecureSkipVerify = true
			tlsConfig.ClientAuth = tls.RequireAnyClientCert
		} else {
			pool := x509.NewCertPool()
			caCrt, err := ioutil.ReadFile(s.option.ClientTrustedCAFile)
			if err != nil {
				return nil, fmt.Errorf("client auth cafile read error: %s", err.Error())
			}

			pool.AppendCertsFromPEM(caCrt)
			tlsConfig.ClientCAs = pool
			tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
		}
	}
	if s.option.ClientCrlFile != "" {
		logger.Infof("ignore client crlfile: %s", s.option.ClientCertFile)
	}

	srv := &http.Server{
		Addr:      s.option.ListenClientURL.Hostname() + ":" + s.option.ListenClientURL.Port(),
		TLSConfig: tlsConfig,
	}

	r := gin.New()
	if s.option.IsPprof {
		pprof.Register(r, nil)
	}

	srv.Handler = r

	ch := make(chan error, 1)
	go func() {
		var err error
		if s.option.IsTLSEnable {
			logger.Infof("Listening and serving HTTPS on %s", srv.Addr)
			err = srv.ListenAndServeTLS(s.option.ClientCertFile, s.option.ClientKeyFile)
		} else {
			logger.Infof("Listening and serving HTTP on %s", srv.Addr)
			err = srv.ListenAndServe()
		}

		if err != nil && err != http.ErrServerClosed {
			logger.Errorf("Server Start Failed: %s", err)
			ch <- err
		}
	}()

	go func() {
		quit := make(chan os.Signal)
		signal.Notify(quit, os.Interrupt)

		var err error
		select {
		case <-quit:
		case err = <-ch:
		}

		// Close on signal
		if err == nil {
			logger.Infof("Shutdown Server ...")

			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()

			if err := srv.Shutdown(ctx); err != nil {
				logger.Errorf("Server Shutdown: %v", err)
			}
		}

		s.stop <- struct{}{}
	}()

	return s.stop, nil
}

// NewServer return Server
func NewServer(option *Options) (Server, error) {
	randomProvider, err := buildRandomProvider(option)
	if err != nil {
		return nil, err
	}
	segmentProvider, err := buildSegmentProvider(option)
	if err != nil {
		return nil, err
	}
	snowflakeProvider, err := buildSnowflakeProvider(option)
	if err != nil {
		return nil, err
	}

	return &server{
		option:            option,
		randomProvider:    randomProvider,
		segmentProvider:   segmentProvider,
		snowflakeProvider: snowflakeProvider,
		stop:              make(chan struct{}),
	}, nil
}

func buildRandomProvider(option *Options) (random.Provider, error) {
	builder, err := randomBuilder.NewBuilder(&randomBuilder.Options{
		ProviderName: option.RandomProvider,
		ProviderArgs: option.RandomProviderArgs,
	})
	if err != nil {
		return nil, err
	}

	return builder.Build()
}

func buildSnowflakeProvider(option *Options) (snowflake.Provider, error) {
	builder, err := snowflakeBuilder.NewBuilder(&snowflakeBuilder.Options{
		ProviderName: option.SnowflakeProvider,
		ProviderArgs: option.SnowflakeProviderArgs,
	})
	if err != nil {
		return nil, err
	}

	return builder.Build()
}

func buildSegmentProvider(option *Options) (segment.Provider, error) {
	builder, err := segmentBuilder.NewBuilder(&segmentBuilder.Options{
		ProviderName: option.SegmentProvider,
		ProviderArgs: option.SegmentProviderArgs,
	})
	if err != nil {
		return nil, err
	}

	return builder.Build()
}
