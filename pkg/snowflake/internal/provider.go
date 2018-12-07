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

// Package internal implement a common snowflake provider, which expect a machine identifier
// to provide machine id.
package internal

import (
	"fmt"
	"sync"

	ierror "github.com/lsytj0413/fyllo/pkg/error"
	"github.com/lsytj0413/fyllo/pkg/snowflake"
)

// Identifier for provice machine id
type Identifier interface {
	// Identify return the machine id
	Identify() (uint64, error)
}

// CommonProvider implement snowflake provider, it want identifier to provide machine id
type CommonProvider struct {
	mutex         sync.RWMutex
	lastMachineID uint64
	providers     map[uint64]snowflake.Provider

	name       string
	identifier Identifier
}

// NewProvider return CommonProvider instance
func NewProvider(name string, identifier Identifier) *CommonProvider {
	return &CommonProvider{
		identifier: identifier,
		name:       name,
		providers:  map[uint64]snowflake.Provider{},
	}
}

// Name implement snowflake.Provider Name
func (p *CommonProvider) Name() string {
	return p.name
}

// Next implement snowflake.Provider Next
func (p *CommonProvider) Next() (*snowflake.Result, error) {
	err := p.checkIdentify()
	if err != nil {
		return nil, err
	}

	provider := p.providerForTag(0)
	r, err := provider.Next()
	if err != nil {
		return nil, err
	}

	r.Name = p.name
	return r, nil
}

func (p *CommonProvider) providerForTag(tag uint64) (provider snowflake.Provider) {
	func() {
		p.mutex.RLock()
		defer p.mutex.RUnlock()

		provider, _ = p.providers[tag]
	}()
	if provider != nil {
		return
	}

	p.mutex.Lock()
	defer p.mutex.Unlock()

	provider = &tagProvider{
		tag:     tag,
		machine: p.lastMachineID,
		t:       newMs(),
		s:       newSequence(snowflake.MaxSequenceValue),
	}
	p.providers[tag] = provider
	return provider
}

// checkIdentify make sure machine id is correct
func (p *CommonProvider) checkIdentify() error {
	currMachineID, err := p.identifier.Identify()
	if err != nil {
		return err
	}

	if currMachineID != p.lastMachineID {
		func() {
			p.mutex.Lock()
			defer p.mutex.Unlock()

			if currMachineID != p.lastMachineID {
				p.lastMachineID = currMachineID
				p.providers = make(map[uint64]snowflake.Provider)
			}
		}()
	}
	return nil
}

type tagProvider struct {
	tag           uint64
	machine       uint64
	lastTimestamp uint64
	mutex         sync.Mutex

	s sequence
	t ms
}

func (p *tagProvider) Name() string {
	return "tagProvider"
}

func (p *tagProvider) Next() (*snowflake.Result, error) {
	_, err := p.processWithTimestamp(p.t.currMicroSeconds())
	if err != nil {
		return nil, err
	}

	sequenceNumber, err := p.s.next()
	if err != nil {
		return nil, err
	}
	r := &snowflake.Result{
		Next: snowflake.MakeSnowflakeID(p.lastTimestamp, p.machine, p.tag, sequenceNumber),
		Labels: map[string]string{
			"timestamp": fmt.Sprintf("%d", p.lastTimestamp),
			"sequence":  fmt.Sprintf("%d", sequenceNumber),
			"tag":       fmt.Sprintf("%d", 0),
			"machine":   fmt.Sprintf("%d", p.machine),
		},
	}
	return r, nil
}

func (p *tagProvider) processWithTimestamp(now uint64) (uint64, error) {
	if now > p.lastTimestamp {
		p.lastTimestamp = now
		p.s.reset()
	} else if now == p.lastTimestamp {
		if p.s.isOutRange() {
			p.lastTimestamp = p.t.waitForNextMS(now)
			p.s.reset()
		}
	} else {
		return 0, ierror.NewError(ierror.EcodeTimestampRewind, fmt.Sprintf("current timestamp[%d] less than last timestamp[%d]", now, p.lastTimestamp))
	}

	return p.lastTimestamp, nil
}
