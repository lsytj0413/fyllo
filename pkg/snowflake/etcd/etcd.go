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
	"context"
	"fmt"
	"strings"

	client "go.etcd.io/etcd/clientv3"

	"github.com/lsytj0413/fyllo/pkg/common"
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
	client *client.Client
}

func (e *etcdIdentifier) Identify() (uint64, error) {
	return e.mid, nil
}

// Register will register the machine id for use
func (e *etcdIdentifier) Register() error {
	prefix := "/fyllo/machine"

	lease, err := e.client.Grant(context.TODO(), 10)
	if err != nil {
		return err
	}
	_, err = e.client.KeepAlive(context.TODO(), lease.ID)
	if err != nil {
		return err
	}

	for i := uint64(0); i < snowflake.MaxMachineValue; i++ {
		key := fmt.Sprintf("%v/%v", prefix, i)
		resp, err := e.client.Txn(context.TODO()).
			If(client.Compare(client.CreateRevision(key), "=", 0)).
			Then(client.OpPut(key, key, client.WithLease(lease.ID))).
			Commit()
		if err != nil {
			return err
		}
		if resp.Succeeded {
			// hold the machine id
			e.mid = i
			return nil
		}
	}
	return fmt.Errorf("machine id register failed, no valid")
}

// Options is etcd snowflake provider option
type Options struct {
	Args string
}

// NewProvider return etcd snowflake provider implement
func NewProvider(options *Options) (snowflake.Provider, error) {
	config, err := parseEtcdConfig(options.Args)
	if err != nil {
		return nil, err
	}

	cli, err := client.New(*config)
	if err != nil {
		return nil, errors.NewError(errors.EcodeInitFailed, fmt.Sprintf("snowflake provider[etcd] init failed, %v", err))
	}

	identifier := &etcdIdentifier{
		client: cli,
	}
	err = identifier.Register()
	if err != nil {
		return nil, errors.NewError(errors.EcodeInitFailed, fmt.Sprintf("snowflake provider[etcd] init failed, %v", err))
	}

	return internal.NewProvider(ProviderName, identifier), nil
}

func parseEtcdConfig(arg string) (*client.Config, error) {
	if 0 == len(arg) {
		return nil, errors.NewError(errors.EcodeInitFailed, fmt.Sprintf("snowflake provider[etcd] argument shoult not be empty"))
	}

	kvs, err := common.SplitKeyValueArrayStringSep(arg, ";")
	if err != nil {
		return nil, errors.NewError(errors.EcodeInitFailed, fmt.Sprintf("snowflake provider[etcd] argument parse failed, %v", err))
	}

	if _, ok := kvs["endpoints"]; !ok {
		return nil, errors.NewError(errors.EcodeInitFailed, fmt.Sprintf("snowflake provider[etcd] argument parse failed, endpoints doesn't exists"))
	}

	config := &client.Config{}
	config.Endpoints = strings.Split(kvs["endpoints"], ",")
	config.Username = kvs["user"]
	config.Password = kvs["pwd"]
	return config, nil
}
