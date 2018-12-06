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

package cmap

import (
	"bytes"
	"sync"
	"sync/atomic"
)

// Bucket is a concurrency safe List of Pair, which has the same hash
type Bucket interface {
	// Put a Pair into Bucket, will store as the linked head
	Put(p Pair, lock sync.Locker) (bool, error)

	// Get the Pair with the given key
	Get(key string) Pair

	// Get the first Pair stored
	GetFirstPair() Pair

	// Delete the given key
	Delete(key string, lock sync.Locker) bool

	// Delete all keys
	Clear(lock sync.Locker)

	// Size returns the count of key
	Size() uint64

	// String returns the readable present of the Bucket
	String() string
}

type bucket struct {
	firstValue atomic.Value
	size       uint64
}

// atomic.Value cannot store nil, this for placeholder
// use size to determine where is have at least one value?
var placeholder Pair = &pair{}

func newBucket() Bucket {
	b := &bucket{}
	b.firstValue.Store(placeholder)
	return b
}

func (b *bucket) GetFirstPair() Pair {
	if v := b.firstValue.Load(); v == nil {
		return nil
	} else if p, ok := v.(Pair); !ok || p == placeholder {
		return nil
	} else {
		return p
	}
}

func (b *bucket) Put(p Pair, lock sync.Locker) (bool, error) {
	if p == nil {
		return false, newInvalidParamError("Bucket.Put: pair is nil")
	}

	if lock != nil {
		lock.Lock()
		defer lock.Unlock()
	}

	firstPair := b.GetFirstPair()
	if firstPair == nil {
		b.firstValue.Store(p)
		atomic.AddUint64(&b.size, 1)
		return true, nil
	}

	// find if need replace
	var target Pair
	key := p.Key()
	for v := firstPair; v != nil; v = v.Next() {
		if v.Key() == key {
			target = v
			break
		}
	}
	// replace the old value
	if target != nil {
		target.SetValue(p.Value())
		return false, nil
	}

	// store new Pair
	p.SetNext(firstPair)
	b.firstValue.Store(p)
	atomic.AddUint64(&b.size, 1)
	return true, nil
}

func (b *bucket) Get(key string) Pair {
	firstPair := b.GetFirstPair()
	if firstPair == nil {
		return nil
	}

	for v := firstPair; v != nil; v = v.Next() {
		if v.Key() == key {
			return v
		}
	}

	return nil
}

func (b *bucket) Delete(key string, lock sync.Locker) bool {
	if lock != nil {
		lock.Lock()
		defer lock.Unlock()
	}

	firstPair := b.GetFirstPair()
	if firstPair == nil {
		return false
	}

	var prevPairs []Pair
	var target Pair
	var breakpoint Pair

	for v := firstPair; v != nil; v = v.Next() {
		if v.Key() == key {
			target = v
			breakpoint = v.Next()
			break
		}
		prevPairs = append(prevPairs, v)
	}
	if target == nil {
		return false
	}

	// link the prevPairs before breakpoint
	// TODO: is this ok?
	newFirstPair := breakpoint
	for i := len(prevPairs) - 1; i >= 0; i-- {
		pairCopy := prevPairs[i].Clone()
		pairCopy.SetNext(newFirstPair)
		newFirstPair = pairCopy
	}
	if newFirstPair != nil {
		b.firstValue.Store(newFirstPair)
	} else {
		b.firstValue.Store(placeholder)
	}

	atomic.AddUint64(&b.size, ^uint64(0))
	return true
}

func (b *bucket) Clear(lock sync.Locker) {
	if lock != nil {
		lock.Lock()
		defer lock.Unlock()
	}

	atomic.StoreUint64(&b.size, 0)
	b.firstValue.Store(placeholder)
}

func (b *bucket) Size() uint64 {
	return atomic.LoadUint64(&b.size)
}

func (b *bucket) String() string {
	var buf bytes.Buffer
	buf.WriteString("[ ")
	for v := b.GetFirstPair(); v != nil; v = v.Next() {
		buf.WriteString(v.String())
		buf.WriteString(" ")
	}

	buf.WriteString("]")
	return buf.String()
}
