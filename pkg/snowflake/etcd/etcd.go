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

// Package etcd implement the provider which's machine id storage by etcd cluster, it only support etcdv3
package etcd

import (
	"fmt"

	"go.etcd.io/etcd/client"

	"github.com/lsytj0413/fyllo/pkg/errors"
	"github.com/lsytj0413/fyllo/pkg/snowflake"
	"github.com/lsytj0413/fyllo/pkg/snowflake/internal"
)

const (
	// ProviderName for the etcd snowflake provider
	ProviderName = "etcd"
)

type etcdIdentifier struct {
	mid    uint64
	client client.Client
}

func (e *etcdIdentifier) Identify() (uint64, error) {
	return e.mid, nil
}

// Options is etcd snowflake provider option
type Options struct {
	Args string
}

// NewProvider return etcd snowflake provider implement
func NewProvider(options *Options) (snowflake.Provider, error) {
	config := client.Config{}
	cli, err := client.New(config)
	if err != nil {
		return nil, errors.NewError(errors.EcodeInitFailed, fmt.Sprintf("snowflake provider[etcd] init failed, %v", err))
	}

	return internal.NewProvider(ProviderName, &etcdIdentifier{
		client: cli,
	}), nil
}
