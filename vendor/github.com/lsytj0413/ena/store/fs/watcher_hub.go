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
	"sync/atomic"
)

// WatcherHub is interface define for Watcher module
type WatcherHub interface {
	watch(key string, recursive bool) (Watcher, error)
	notify(*Result)
}

// watcherHub contains all subscribed watchers
type watcherHub struct {
	sync.Mutex

	count    int64 // current number of watchers
	watchers map[string]*list.List
}

func newWatchHub(capacity int) WatcherHub {
	return &watcherHub{
		watchers: make(map[string]*list.List),
	}
}

// watch function returns a Watcher
// if recursive is true, the first change under key will be sent to the event channel of the watcher
// if recursive is false, the first change at key will be sent to the event channel of the watcher
func (h *watcherHub) watch(key string, recursive bool) (Watcher, error) {
	w := &defWatcher{
		resultChan: make(chan *Result, 100),
		recursive:  recursive,
		stream:     true,
		hub:        h,
	}
	h.Lock()
	defer h.Unlock()

	l, ok := h.watchers[key]
	if !ok {
		l = list.New()
		h.watchers[key] = l
	}
	elem := l.PushBack(w)

	w.remove = func() {
		if w.removed {
			return
		}

		w.removed = true
		l.Remove(elem)
		atomic.AddInt64(&h.count, -1)
		if l.Len() == 0 {
			delete(h.watchers, key)
		}
	}
	atomic.AddInt64(&h.count, 1)
	return w, nil
}

func (h *watcherHub) add(r *Result) {
	// h.resultHistory.addResult(r)
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
	h.Lock()
	defer h.Unlock()

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
