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

// Package notify enables independent components of an application to
// observe notable events in a decoupled fashion.
// It generalized the pattern of *multiple* consumers of an event.
package notify

import (
	"errors"
	"sync"
	"time"
)

var (
	// ErrEventNotFound is error for event not exists
	ErrEventNotFound = errors.New("ErrEventNotFound")
)

// Notifier is interface define for Notify
type Notifier interface {
	On(event string, ch chan interface{}) error
	Remove(event string, ch chan interface{}) error
	RemoveAll(event string) error
	Emit(event string, data interface{}) error
	EmitWithTimeout(event string, data interface{}, timeout time.Duration) error
}

var (
	defaultNotifier Notifier
)

func init() {
	defaultNotifier = New()
}

// Default returns the default notifier
func Default() Notifier {
	return defaultNotifier
}

// New returns a new notifier
func New() Notifier {
	return &defNotifier{
		events: make(map[string][]chan interface{}),
	}
}

type defNotifier struct {
	events map[string][]chan interface{}
	mutex  sync.RWMutex
}

func (n *defNotifier) On(event string, ch chan interface{}) error {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	n.events[event] = append(n.events[event], ch)
	return nil
}

func (n *defNotifier) Remove(event string, ch chan interface{}) error {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	newArray := make([]chan interface{}, 0)
	outChans, ok := n.events[event]
	if !ok {
		return ErrEventNotFound
	}

	for _, outCh := range outChans {
		if outCh != ch {
			newArray = append(newArray, ch)
		} else {
			close(outCh)
		}
	}
	n.events[event] = newArray
	return nil
}

func (n *defNotifier) RemoveAll(event string) error {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	outChans, ok := n.events[event]
	if !ok {
		return ErrEventNotFound
	}
	for _, outCh := range outChans {
		close(outCh)
	}
	delete(n.events, event)
	return nil
}

func (n *defNotifier) Emit(event string, data interface{}) error {
	n.mutex.RLock()
	defer n.mutex.RUnlock()

	outChans, ok := n.events[event]
	if !ok {
		return ErrEventNotFound
	}
	for _, outCh := range outChans {
		outCh <- data
	}
	return nil
}

func (n *defNotifier) EmitWithTimeout(event string, data interface{}, timeout time.Duration) error {
	n.mutex.RLock()
	defer n.mutex.RUnlock()

	outChans, ok := n.events[event]
	if !ok {
		return ErrEventNotFound
	}
	for _, outCh := range outChans {
		select {
		case outCh <- data:
		case <-time.After(timeout):
		}
	}
	return nil
}
