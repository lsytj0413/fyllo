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

import "sync"

// ResultHistory is Action history store
type ResultHistory struct {
	Queue *resultQueue
	rwl   sync.RWMutex
}

func newResultHistory(capacity int) *ResultHistory {
	return &ResultHistory{
		Queue: &resultQueue{
			Capacity: capacity,
			Results:  make([]*Result, capacity),
		},
	}
}

func (h *ResultHistory) addResult(r *Result) *Result {
	h.rwl.Lock()
	defer h.rwl.Unlock()

	h.Queue.insert(r)
	return r
}

func (h *ResultHistory) clone() *ResultHistory {
	return &ResultHistory{
		Queue: h.Queue.clone(),
	}
}
