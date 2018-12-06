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

package conc

import "testing"

func Benchmark_Set(b *testing.B) {
	arr, _ := NewConcurrentArray(100)

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			for j := uint32(0); j < arr.Len(); j++ {
				arr.Set(j, nil)
			}
		}
	})
	// for i := 0; i < b.N; i++ {
	// 	for j := uint32(0); j < arr.Len(); j++ {
	// 		arr.Set(j, nil)
	// 	}
	// }
}

func Benchmark_Get(b *testing.B) {
	arr, _ := NewConcurrentArray(100)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for j := uint32(0); j < arr.Len(); j++ {
				arr.Get(j)
			}
		}
	})

	// for i := 0; i < b.N; i++ {
	// 	for j := uint32(0); j < arr.Len(); j++ {
	// 		arr.Get(j)
	// 	}
	// }
}

func Benchmark_SetGet(b *testing.B) {
	arr, _ := NewConcurrentArray(100)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for j := uint32(0); j < arr.Len(); j++ {
				arr.Set(j, nil)
				arr.Get(j)
			}
		}
	})
	// for i := 0; i < b.N; i++ {
	// 	for j := uint32(0); j < arr.Len(); j++ {
	// 		arr.Set(j, nil)
	// 		arr.Get(j)
	// 	}
	// }
}
