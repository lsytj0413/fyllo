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

package internal

import (
	"sync"
)

// Result for Next
type Result struct {
	Next      uint64 `json:"next"`
	Timestamp uint64 `json:"timestamp"`
	Sequence  uint64 `json:"sequence"`
	Machine   uint64 `json:"machine"`
	Tag       uint64 `json:"tag"`
}

// Generator interface for snowflake generator
type Generator interface {
	Next(uint64) (*Result, error)
}

// default Generator implement
type sfGenerator struct {
	mutex   sync.RWMutex
	machine uint64
	mappers map[uint64]Generator
}

func (sf *sfGenerator) getGeneratorByTag(tag uint64) Generator {
	getter := func() Generator {
		sf.mutex.RLock()
		defer sf.mutex.RUnlock()

		g, ok := sf.mappers[tag]
		if !ok {
			return nil
		}

		return g
	}

	g := getter()
	if g != nil {
		return g
	}

	sf.mutex.Lock()
	defer sf.mutex.Unlock()

	// try get again
	g, ok := sf.mappers[tag]
	if ok {
		return g
	}

	g = &tagGenerator{
		tag:     tag,
		machine: sf.machine,
	}
	sf.mappers[tag] = g
	return g
}

func (sf *sfGenerator) Next(tag uint64) (*Result, error) {
	g := sf.getGeneratorByTag(tag)
	return g.Next(tag)
}

// NewGenerator construct a Generator instance
func NewGenerator(machine uint64) (Generator, error) {
	g := &sfGenerator{
		mutex:   sync.RWMutex{},
		machine: machine,
		mappers: make(map[uint64]Generator, 32),
	}

	return g, nil
}
