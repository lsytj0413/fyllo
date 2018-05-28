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

package segment

import (
	"context"

	"github.com/lsytj0413/fyllo/conf"
)

type generator struct {
	ch <-chan uint64
}

func (g *generator) Next(c context.Context, tag string) (*conf.SegmentResult, error) {
	r := &conf.SegmentResult{
		Tag: tag,
	}
	select {
	case id := <-g.ch:
		r.Next = id
		return r, nil
	case <-c.Done():
		return nil, c.Err()
	}
}
