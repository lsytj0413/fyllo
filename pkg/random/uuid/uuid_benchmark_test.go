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

package uuid

import "testing"

func Benchmark_Next_Parallel(b *testing.B) {
	p, err := NewProvider(nil)
	if err != nil {
		b.Fatalf("NewProvider failed, %v", err)
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err = p.Next(nil)
			if err != nil {
				b.Fatalf("Next failed, %v", err)
			}
		}
	})
}

func Benchmark_Next(b *testing.B) {
	p, err := NewProvider(nil)
	if err != nil {
		b.Fatalf("NewProvider failed, %v", err)
	}

	for i := 0; i < b.N; i++ {
		_, err = p.Next(nil)
		if err != nil {
			b.Fatalf("Next failed, %v", err)
		}
	}
}
