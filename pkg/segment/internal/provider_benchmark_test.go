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

package internal

import (
	"math/rand"
	"testing"

	"github.com/lsytj0413/fyllo/pkg/segment"
)

func Benchmark_Next_Parallel(b *testing.B) {
	tagItems := []*TagItem{
		{
			Tag:         "1",
			Max:         10,
			Min:         1,
			Description: "1",
		},
		{
			Tag:         "2",
			Max:         20,
			Min:         2,
			Description: "2",
		},
	}
	mockStorager := newMockStorager(tagItems)

	providerName := "test"
	p, err := NewProvider(providerName, mockStorager)

	if err != nil {
		b.Fatalf("NewProvider failed, %v", err)
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			index := rand.Intn(len(tagItems))
			_, err = p.Next(&segment.Arguments{
				Tag: tagItems[index].Tag,
			})
			if err != nil {
				b.Fatalf("Next failed, %v", err)
			}
		}
	})
}

func Benchmark_Next(b *testing.B) {
	tagItems := []*TagItem{
		{
			Tag:         "1",
			Max:         10,
			Min:         1,
			Description: "1",
		},
		{
			Tag:         "2",
			Max:         20,
			Min:         2,
			Description: "2",
		},
	}
	mockStorager := newMockStorager(tagItems)

	providerName := "test"
	p, err := NewProvider(providerName, mockStorager)

	if err != nil {
		b.Fatalf("NewProvider failed, %v", err)
	}

	for i := 0; i < b.N; i++ {
		index := rand.Intn(len(tagItems))
		_, err = p.Next(&segment.Arguments{
			Tag: tagItems[index].Tag,
		})
		if err != nil {
			b.Fatalf("Next failed, %v", err)
		}
	}
}
