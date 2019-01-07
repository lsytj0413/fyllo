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

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/lsytj0413/ena/logger"
)

// Server is fyllo server interface
type Server interface {
	Options() *Options
	Start(api Installer, stop chan struct{}) error
	Shutdown(ctx context.Context) error
}

// Installer for install api
type Installer interface {
	Install(engine *gin.Engine) error
}

type funcInstaller struct {
	fn func(*gin.Engine) error
}

func (i *funcInstaller) Install(engine *gin.Engine) error {
	return i.fn(engine)
}

// NewInstallerFunction return Installer from function
func NewInstallerFunction(f func(*gin.Engine) error) Installer {
	return &funcInstaller{
		fn: f,
	}
}

// NewInstallers combain multi installer
func NewInstallers(installers ...Installer) Installer {
	return NewInstallerFunction(func(engine *gin.Engine) (err error) {
		for i := 0; i < len(installers); i++ {
			err = installers[i].Install(engine)
			if err != nil {
				return
			}
		}
		return
	})
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
	IsPprof                bool     `json:"isPprof"`

	IsTLSEnable bool `json:"isTLSEnable"` // auto filled by ListenClientURL
}

type server struct {
	option *Options
	srv    *http.Server
}

func (s *server) Options() *Options {
	return s.option
}

func (s *server) Start(api Installer, stop chan struct{}) error {
	if s.option.ListenClientURL.Scheme == "https" {
		if s.option.ClientCertFile == "" || s.option.ClientKeyFile == "" {
			return fmt.Errorf("listen on https without keyfile or certfile exists")
		}

		s.option.IsTLSEnable = true
	}

	if s.option.IsClientCertAuthEnable {
		if s.option.ClientTrustedCAFile == "" && !s.option.IsInsecureSkipVerify {
			return fmt.Errorf("client auth enable without client-trusted-ca-file or client-auto-tls set")
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
				return fmt.Errorf("client auth cafile read error: %s", err.Error())
			}

			pool.AppendCertsFromPEM(caCrt)
			tlsConfig.ClientCAs = pool
			tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
		}
	}
	if s.option.ClientCrlFile != "" {
		logger.Infof("ignore client crlfile: %s", s.option.ClientCertFile)
	}

	s.srv = &http.Server{
		Addr:      s.option.ListenClientURL.Hostname() + ":" + s.option.ListenClientURL.Port(),
		TLSConfig: tlsConfig,
	}

	r := gin.New()
	if s.option.IsPprof {
		pprof.Register(r, nil)
	}
	err := api.Install(r)
	if err != nil {
		return err
	}

	s.srv.Handler = r

	if s.option.IsTLSEnable {
		logger.Infof("Listening and serving HTTPS on %s", s.srv.Addr)
		err = s.srv.ListenAndServeTLS(s.option.ClientCertFile, s.option.ClientKeyFile)
	} else {
		logger.Infof("Listening and serving HTTP on %s", s.srv.Addr)
		err = s.srv.ListenAndServe()
	}

	stop <- struct{}{}

	if err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

// NewServer return Server
func NewServer(option *Options) (Server, error) {
	if option.ListenClientURL.Scheme == "https" {
		if option.ClientCertFile == "" || option.ClientKeyFile == "" {
			return nil, fmt.Errorf("listen on https without keyfile or certfile exists")
		}

		option.IsTLSEnable = true
	}

	if option.IsClientCertAuthEnable {
		if option.ClientTrustedCAFile == "" && !option.IsInsecureSkipVerify {
			return nil, fmt.Errorf("client auth enable without client-trusted-ca-file or client-auto-tls set")
		}
	}

	return &server{
		option: option,
	}, nil
}
