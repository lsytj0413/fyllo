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
	"fmt"
	"sync"
	"sync/atomic"
)

// Segment is concurrency safe defines, which store Bucket and has a lock
type Segment interface {
	Put(p Pair) (bool, error)
	Get(key string) Pair
	GetWithHash(key string, keyHash uint64) Pair
	Delete(key string) bool
	Size() uint64
}

type segment struct {
	buckets           []Bucket
	bucketLen         uint32
	pairTotal         uint64
	pairRedistributor PairRedistributor
	lock              sync.Mutex
}

func newSegment(bucketNumber uint32, pairRedistributor PairRedistributor) Segment {
	if bucketNumber < 0 {
		bucketNumber = DefaultBucketNumber
	}

	if pairRedistributor == nil {
		pairRedistributor = newPairRedistributor(DefaultBucketLoadFactor, bucketNumber)
	}

	buckets := make([]Bucket, bucketNumber)
	for i := uint32(0); i < bucketNumber; i++ {
		buckets[i] = newBucket()
	}

	return &segment{
		buckets:           buckets,
		bucketLen:         bucketNumber,
		pairRedistributor: pairRedistributor,
	}
}

func (s *segment) Put(p Pair) (bool, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	b := s.buckets[int(p.Hash()%uint64(s.bucketLen))]
	ok, err := b.Put(p, nil)
	if ok {
		newTotal := atomic.AddUint64(&s.pairTotal, 1)
		s.redistribute(newTotal, b.Size())
	}

	return ok, err
}

func (s *segment) Get(key string) Pair {
	return s.GetWithHash(key, hash(key))
}

func (s *segment) GetWithHash(key string, keyHash uint64) Pair {
	b := func() Bucket {
		s.lock.Lock()
		defer s.lock.Unlock()

		return s.buckets[int(keyHash%uint64(s.bucketLen))]
	}()

	return b.Get(key)
}

func (s *segment) Delete(key string) bool {
	s.lock.Lock()
	defer s.lock.Unlock()

	b := s.buckets[int(hash(key)%uint64(s.bucketLen))]
	ok := b.Delete(key, nil)
	if ok {
		newTotal := atomic.AddUint64(&s.pairTotal, ^uint64(0))
		s.redistribute(newTotal, b.Size())
	}

	return ok
}

func (s *segment) Size() uint64 {
	return atomic.LoadUint64(&s.pairTotal)
}

func (s *segment) redistribute(pairTotal uint64, bucketSize uint64) (err error) {
	defer func() {
		if p := recover(); p != nil {
			if pErr, ok := p.(error); ok {
				err = newPairRedistributorError(pErr.Error())
			} else {
				err = newPairRedistributorError(fmt.Sprintf("%s", p))
			}
		}
	}()

	s.pairRedistributor.UpdateThreshold(pairTotal, s.bucketLen)
	bucketStatus := s.pairRedistributor.CheckBucketStatus(pairTotal, bucketSize)
	newBuckets, changed := s.pairRedistributor.Redistribe(bucketStatus, s.buckets)
	if changed {
		s.buckets = newBuckets
		s.bucketLen = uint32(len(s.buckets))
	}

	return nil
}
