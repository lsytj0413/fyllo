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
	"sync"

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
	seqs          map[uint64]Sequencer

	name       string
	identifier Identifier
}

// NewProvider return CommonProvider instance
func NewProvider(name string, identifier Identifier) *CommonProvider {
	return &CommonProvider{
		identifier: identifier,
		name:       name,
		seqs:       map[uint64]Sequencer{},
	}
}

// Name implement snowflake.Provider Name
func (p *CommonProvider) Name() string {
	return p.name
}

// Next implement snowflake.Provider Next
func (p *CommonProvider) Next(arg *snowflake.Arguments) (*snowflake.Result, error) {
	err := p.checkIdentify()
	if err != nil {
		return nil, err
	}

	seqer := p.sequencerForTag(0)
	sequenceNumber, timestamp, err := seqer.Next()
	if err != nil {
		return nil, err
	}

	r := &snowflake.Result{
		Name: p.name,
		Next: snowflake.MakeSnowflakeID(timestamp, p.lastMachineID, 0, sequenceNumber),
		Labels: map[string]string{
			"sequence":  "",
			"timestamp": "",
			"tag":       "0",
			"machine":   "",
		},
	}
	return r, nil
}

func (p *CommonProvider) sequencerForTag(tag uint64) (seq Sequencer) {
	func() {
		p.mutex.RLock()
		defer p.mutex.RUnlock()

		seq, _ = p.seqs[tag]
	}()
	if seq != nil {
		return
	}

	p.mutex.Lock()
	defer p.mutex.Unlock()

	seq = NewSequencer()
	p.seqs[tag] = seq
	return seq
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
				p.seqs = make(map[uint64]Sequencer)
			}
		}()
	}
	return nil
}
