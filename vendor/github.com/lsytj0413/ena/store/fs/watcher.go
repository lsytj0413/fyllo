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

package fs

// Watcher defines a interface for notice the store operate Result
type Watcher interface {
	ResultChan() <-chan *Result
	Remove()
	notify(*Result, bool, bool) bool
}

type defWatcher struct {
	resultChan chan *Result
	stream     bool
	recursive  bool

	hub     *watcherHub
	removed bool
	remove  func()
}

func (w *defWatcher) ResultChan() <-chan *Result {
	return w.resultChan
}

func (w *defWatcher) Remove() {
	w.hub.Lock()
	defer w.hub.Unlock()

	close(w.resultChan)
	if w.remove != nil {
		w.remove()
	}
}

// notify function notifies the watcher. If the watcher interests in the given path,
// the function will return true
func (w *defWatcher) notify(r *Result, originalPath bool, deleted bool) bool {
	// watcher is interested the path in three cases and under one condition
	// the condition is that the event happens after the warcher's sinceIndex

	// 1. the path at which the event happens is the path the watcher is watching at.
	// 2. the watcher is a recursive watcher, it interested in the event happens after its watching path.
	// 3. when we delete a directory, we need to force notify all the watchers who watches at the file we need to delete.
	if w.recursive || originalPath || deleted {
		select {
		case w.resultChan <- r:
		default:
			// missed, remove the watcher
			w.Remove()
		}
		return true
	}

	return false
}
