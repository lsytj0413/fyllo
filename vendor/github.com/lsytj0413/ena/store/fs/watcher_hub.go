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

import (
	"container/list"
	"path"
	"strings"
	"sync"
)

// WatcherHub is interface define for Watcher module
type WatcherHub interface {
	watch(key string, recursive bool) (Watcher, error)
	notify(*Result)
}

// watcherHub contains all subscribed watchers
type watcherHub struct {
	count uint64 // current number of watchers

	mutex         sync.Mutex
	watchers      map[string]*list.List
	resultHistory *ResultHistory
}

func newWatchHub(capacity int) WatcherHub {
	return &watcherHub{
		watchers:      make(map[string]*list.List),
		resultHistory: newResultHistory(capacity),
	}
}

func (h *watcherHub) watch(key string, recursive bool) (Watcher, error) {
	return nil, nil
}

func (h *watcherHub) add(r *Result) {
	h.resultHistory.addResult(r)
}

func (h *watcherHub) notify(r *Result) {
	h.add(r)

	segments := components(r.CurrNode.Key)
	currPath := "/"
	for _, segment := range segments {
		currPath = keyFromDirAndFile(currPath, segment)
		h.notifyWatchers(r, currPath, false)
	}
}

func (h *watcherHub) notifyWatchers(r *Result, nodePath string, deleted bool) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	watchers, ok := h.watchers[nodePath]
	if ok {
		curr := watchers.Front()
		for curr != nil {
			next := curr.Next()
			w, _ := curr.Value.(Watcher)

			originalPath := (r.CurrNode.Key == nodePath)
			if (originalPath || !isHidden(nodePath, r.CurrNode.Key)) &&
				w.notify(r, originalPath, deleted) {

			}

			curr = next
		}

		if watchers.Len() == 0 {
			delete(h.watchers, nodePath)
		}
	}
}

func isHidden(watchPath string, keyPath string) bool {
	if len(watchPath) > len(keyPath) {
		return false
	}

	afterPath := path.Clean("/" + keyPath[len(watchPath):])
	return strings.Contains(afterPath, "/.")
}
