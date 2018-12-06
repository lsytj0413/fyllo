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

import "sync/atomic"

// BucketStatus defines Bucket LoadFactor status
type BucketStatus uint8

const (
	// BucketStatusNormal present normal LoadFactor
	BucketStatusNormal BucketStatus = 0

	// BucketStatusUnderWeight present under weight LoadFactor
	BucketStatusUnderWeight BucketStatus = 1

	// BucketStatusOverWeight present over weight LoadFactor
	BucketStatusOverWeight BucketStatus = 2
)

// PairRedistributor defines for pair redistributor in Bucket
type PairRedistributor interface {
	UpdateThreshold(pairTotal uint64, bucketNumber uint32)
	CheckBucketStatus(pairTotal uint64, bucketSize uint64) BucketStatus
	Redistribe(bucketStatus BucketStatus, buckets []Bucket) ([]Bucket, bool)
}

type defPairRedistributor struct {
	loadFactor            float64
	upperThreshold        uint64
	overWeightBucketCount uint64
	emptyBucketCount      uint64
}

func newPairRedistributor(loadFactor float64, bucketNumber uint32) PairRedistributor {
	if loadFactor < 0 {
		loadFactor = DefaultBucketLoadFactor
	}

	pr := &defPairRedistributor{}
	pr.loadFactor = loadFactor
	pr.UpdateThreshold(0, bucketNumber)
	return pr
}

var bucketCountTemplate = `Bucket count:
    pairTotal: %d
    bucketNumber: %d
    average: %f
    upperThreshold: %d
    emptyBucketCount: %d

`

func (p *defPairRedistributor) UpdateThreshold(pairTotal uint64, bucketNumber uint32) {
	var average float64
	average = float64(pairTotal / uint64(bucketNumber))
	if average < 100 {
		average = 100
	}

	atomic.StoreUint64(&p.upperThreshold, uint64(average*p.loadFactor))
}

var bucketStatusTemplate = `Check bucket status:
    pairTotal: %d
    bucketSize: %d
    upperThreshold: %d
    overWeightBucketCount: %d
    emptyBucketCount: %d
    bucketStatus: %d

`

func (p *defPairRedistributor) CheckBucketStatus(pairTotal uint64, bucketSize uint64) BucketStatus {
	if bucketSize > DefaultBucketMaxSize ||
		bucketSize >= atomic.LoadUint64(&p.upperThreshold) {
		atomic.AddUint64(&p.overWeightBucketCount, 1)
		return BucketStatusOverWeight
	}

	if bucketSize == 0 {
		atomic.AddUint64(&p.emptyBucketCount, 1)
	}
	return BucketStatusNormal
}

var redistributionTemplate = `Redistributing:
    bucketStatus: %d
    currentNumber: %d
    newNumber: %d

`

func (p *defPairRedistributor) Redistribe(bucketStatus BucketStatus, buckets []Bucket) (newBuckets []Bucket, changed bool) {
	currentNumber := uint64(len(buckets))
	newNumber := currentNumber

	switch bucketStatus {
	case BucketStatusOverWeight:
		if atomic.LoadUint64(&p.overWeightBucketCount)*4 < currentNumber {
			return nil, false
		}
		newNumber = currentNumber << 1
	case BucketStatusUnderWeight:
		if currentNumber < 100 ||
			atomic.LoadUint64(&p.emptyBucketCount)*4 < currentNumber {
			return nil, false
		}
		newNumber = currentNumber >> 1
		if newNumber < 2 {
			newNumber = 2
		}
	default:
		return nil, false
	}

	if newNumber == currentNumber {
		atomic.StoreUint64(&p.overWeightBucketCount, 0)
		atomic.StoreUint64(&p.emptyBucketCount, 0)
		return nil, false
	}

	var pairs []Pair
	for _, b := range buckets {
		for e := b.GetFirstPair(); e != nil; e = e.Next() {
			pairs = append(pairs, e)
		}
	}

	if newNumber > currentNumber {
		for i := uint64(0); i < currentNumber; i++ {
			buckets[i].Clear(nil)
		}
		for j := newNumber - currentNumber; j > 0; j-- {
			buckets = append(buckets, newBucket())
		}
	} else {
		buckets = make([]Bucket, newNumber)
		for i := uint64(0); i < newNumber; i++ {
			buckets[i] = newBucket()
		}
	}

	var count int
	for _, p := range pairs {
		index := int(p.Hash() % newNumber)
		b := buckets[index]
		b.Put(p, nil)
		count++
	}

	atomic.StoreUint64(&p.overWeightBucketCount, 0)
	atomic.StoreUint64(&p.emptyBucketCount, 0)
	return buckets, true
}
