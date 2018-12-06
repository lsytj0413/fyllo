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

type resultQueue struct {
	Results  []*Result
	Size     int
	Front    int
	Back     int
	Capacity int
}

func (q *resultQueue) insert(r *Result) {
	q.Results[q.Back] = r
	q.Back = (q.Back + 1) % q.Capacity

	if q.Size == q.Capacity {
		q.Front = (q.Front + 1) % q.Capacity
	} else {
		q.Size++
	}
}

func (q *resultQueue) clone() *resultQueue {
	cloned := &resultQueue{
		Results:  make([]*Result, q.Capacity),
		Size:     q.Size,
		Front:    q.Front,
		Back:     q.Back,
		Capacity: q.Capacity,
	}

	for i, r := range q.Results {
		cloned.Results[i] = r
	}

	return cloned
}
