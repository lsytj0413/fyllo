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

// Package wait provides utility functions for polling, listening using Go channel
package wait

import (
	"fmt"
	"sync"
)

// Wait is an interface that provices the ability to wait and trigger events that
// are associated with ID
type Wait interface {
	// Register waits returns a chan that waits on the given ID.
	// The chan will be triggered when Trigger is called with the same ID
	Register(id uint64) (<-chan interface{}, error)
	// Trigger triggers the waiting chans with the given ID
	Trigger(id uint64, x interface{}) error
	IsRegistered(id uint64) bool
}

type defWait struct {
	mutex sync.RWMutex
	m     map[uint64]chan interface{}
}

const (
	defaultMapSize = 1000
)

// New creates a Wait object
func New() Wait {
	return &defWait{
		m: make(map[uint64]chan interface{}, defaultMapSize),
	}
}

func (w *defWait) Register(id uint64) (<-chan interface{}, error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if _, ok := w.m[id]; !ok {
		c := make(chan interface{}, 1)
		w.m[id] = c
		return c, nil
	}

	return nil, fmt.Errorf("Wait.Register Failed: duplicate id=%x", id)
}

func (w *defWait) Trigger(id uint64, x interface{}) error {
	f := func() (chan interface{}, bool) {
		w.mutex.Lock()
		defer w.mutex.Unlock()

		c, ok := w.m[id]
		delete(w.m, id)
		return c, ok
	}

	c, ok := f()
	if !ok {
		return fmt.Errorf("Wait.Trigger Failed: id=%d not registered", id)
	}

	select {
	case c <- x:
		close(c)
		return nil
	default:
		return fmt.Errorf("Wait.Trigger Failed: id=%d timeout", id)
	}
}

func (w *defWait) IsRegistered(id uint64) bool {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	_, ok := w.m[id]
	return ok
}
