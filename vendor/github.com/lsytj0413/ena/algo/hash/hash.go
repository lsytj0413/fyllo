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

package hash

import (
	"hash"
	"hash/fnv"
	"sync"
)

// Hasher is interface for Hash
type Hasher interface {
	Uint64(data []byte) uint64
	Uint32(data []byte) uint32
}

type hasher struct {
	v32 hash.Hash32
	v64 hash.Hash64
}

func (h *hasher) Uint64(data []byte) uint64 {
	defer h.v64.Reset()

	h.v64.Write(data)
	return h.v64.Sum64()
}

func (h *hasher) Uint32(data []byte) uint32 {
	defer h.v32.Reset()

	h.v32.Write(data)
	return h.v32.Sum32()
}

// NewHash returns Hasher implement, not thread safe
func NewHash() Hasher {
	return &hasher{
		v32: fnv.New32(),
		v64: fnv.New64(),
	}
}

type safeHasher struct {
	Hasher
	sync.Mutex
}

func (h *safeHasher) Uint32(data []byte) uint32 {
	h.Lock()
	defer h.Unlock()

	return h.Hasher.Uint32(data)
}

func (h *safeHasher) Uint64(data []byte) uint64 {
	h.Lock()
	defer h.Unlock()

	return h.Hasher.Uint64(data)
}

// NewSafeHash returns Hasher implement, thread safe
func NewSafeHash() Hasher {
	return &safeHasher{
		Hasher: NewHash(),
	}
}

var (
	defHasher     = NewHash()
	defSafeHasher = NewSafeHash()
)

// Uint32 return hash uint32
func Uint32(data []byte) uint32 {
	return defHasher.Uint32(data)
}

// Uint64 return hash uint64
func Uint64(data []byte) uint64 {
	return defHasher.Uint64(data)
}

// SafeUint32 return hash uint32, thread safe
func SafeUint32(data []byte) uint32 {
	return defSafeHasher.Uint32(data)
}

// SafeUint64 return hash uint64, thread safe
func SafeUint64(data []byte) uint64 {
	return defSafeHasher.Uint64(data)
}
