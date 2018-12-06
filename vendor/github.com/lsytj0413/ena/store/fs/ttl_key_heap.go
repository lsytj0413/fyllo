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

import "container/heap"

// keyTTLHeap is a min-heap of TTL Keys order by expiration time
type keyTTLHeap interface {
}

// defaultKeyHeap must implements heap.Interface(include sort.Interface)
type defaultKeyHeap struct {
	array []*inode
	// keep the array.index for each inode
	keyMap map[*inode]int
}

func newKeyTTLHeap() keyTTLHeap {
	h := &defaultKeyHeap{keyMap: make(map[*inode]int)}
	heap.Init(h)
	return h
}

// Implement sort.Interface.Len
func (h *defaultKeyHeap) Len() int {
	return len(h.array)
}

// Implement sort.Interface.Less
func (h *defaultKeyHeap) Less(i int, j int) bool {
	return h.array[i].Expire.Before(h.array[j].Expire)
}

// Implement sort.Interface.Swap
func (h *defaultKeyHeap) Swap(i int, j int) {
	// swap inode in array
	h.array[i], h.array[j] = h.array[j], h.array[i]

	// update keyMap
	h.keyMap[h.array[i]] = i
	h.keyMap[h.array[j]] = j
}

// Implement heap.Interface.Push
func (h *defaultKeyHeap) Push(x interface{}) {
	n, _ := x.(*inode)
	h.keyMap[n] = len(h.array)
	h.array = append(h.array, n)
}

// Implement heap.Interface.Pop
func (h *defaultKeyHeap) Pop() interface{} {
	old := h.array
	n := len(old)
	x := old[n-1]

	// set element to nil for GC
	old[n-1] = nil
	h.array = old[0 : n-1]
	delete(h.keyMap, x)
	return x
}

func (h *defaultKeyHeap) top() *inode {
	if h.Len() != 0 {
		return h.array[0]
	}

	return nil
}

func (h *defaultKeyHeap) pop() *inode {
	x := heap.Pop(h)
	n, _ := x.(*inode)
	return n
}

func (h *defaultKeyHeap) push(x interface{}) {
	heap.Push(h, x)
}

func (h *defaultKeyHeap) update(n *inode) {
	index, ok := h.keyMap[n]
	if ok {
		heap.Remove(h, index)
		heap.Push(h, n)
	}
}

func (h *defaultKeyHeap) remove(n *inode) {
	index, ok := h.keyMap[n]
	if ok {
		heap.Remove(h, index)
	}
}
