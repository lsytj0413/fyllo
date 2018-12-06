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

const (
	// DefaultBucketLoadFactor is the default LoadFactor,
	// when LoadFactor greater than this value, the map will re-hash
	DefaultBucketLoadFactor float64 = 0.75

	// DefaultBucketNumber is the default bucket number in a hash segment
	DefaultBucketNumber uint32 = 16

	// DefaultBucketMaxSize is the max size of bucket
	DefaultBucketMaxSize uint64 = 1000

	// MaxConcurrency is the max concurrency number
	MaxConcurrency uint32 = 65536
)
