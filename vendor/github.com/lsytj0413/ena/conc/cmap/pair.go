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
	"fmt"
	"sync/atomic"
	"unsafe"
)

// Pair is a define for key-value item
type Pair interface {
	linkedPair
	Key() string
	Hash() uint64
	Value() interface{}
	SetValue(v interface{}) error
	Clone() Pair
	String() string
}

type pair struct {
	key   string
	hash  uint64
	value unsafe.Pointer
	next  unsafe.Pointer
}

type linkedPair interface {
	Next() Pair
	SetNext(n Pair) error
}

func newPair(key string, v interface{}) (Pair, error) {
	if v == nil {
		return nil, newInvalidParamError("newPair: value is nil")
	}

	p := &pair{
		key:  key,
		hash: hash(key),
	}

	p.value = unsafe.Pointer(&v)
	return p, nil
}

func (p *pair) Key() string {
	return p.key
}

func (p *pair) Hash() uint64 {
	return p.hash
}

func (p *pair) Value() interface{} {
	pointer := atomic.LoadPointer(&p.value)
	if pointer == nil {
		return nil
	}

	return *(*interface{})(pointer)
}

func (p *pair) SetValue(v interface{}) error {
	if v == nil {
		return newInvalidParamError("Pair.SetValue: value is nil")
	}

	atomic.StorePointer(&p.value, unsafe.Pointer(&v))
	return nil
}

func (p *pair) Next() Pair {
	pointer := atomic.LoadPointer(&p.next)
	if pointer == nil {
		return nil
	}

	return (*pair)(pointer)
}

func (p *pair) SetNext(n Pair) error {
	if n == nil {
		atomic.StorePointer(&p.next, nil)
		return nil
	}

	pp, ok := n.(*pair)
	if !ok {
		return newInvalidPairTypeError(fmt.Sprintf("%T", n))
	}

	atomic.StorePointer(&p.next, unsafe.Pointer(pp))
	return nil
}

func (p *pair) Clone() Pair {
	c, _ := newPair(p.Key(), p.Value())
	return c
}

func (p *pair) String() string {
	return p.toString(false)
}

func (p *pair) toString(withNext bool) string {
	var buf bytes.Buffer
	buf.WriteString("pair{key:")
	buf.WriteString(p.Key())
	buf.WriteString(", hash:")
	buf.WriteString(fmt.Sprintf("%d", p.Hash()))
	buf.WriteString(", element:")
	buf.WriteString(fmt.Sprintf("%+v", p.Value()))
	if withNext {
		buf.WriteString(", next:")
		if next := p.Next(); next != nil {
			if npp, ok := next.(*pair); ok {
				buf.WriteString(npp.toString(withNext))
			} else {
				buf.WriteString("<ignore>")
			}
		}
	} else {
		buf.WriteString(", nextKey:")
		if next := p.Next(); next != nil {
			buf.WriteString(next.Key())
		}
	}

	buf.WriteString("}")
	return buf.String()
}
