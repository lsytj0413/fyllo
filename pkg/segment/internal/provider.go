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
	"fmt"
	"sync"

	"github.com/lsytj0413/fyllo/pkg/segment"
)

// Storager is storage for segment provider
type Storager interface {
	// List return all available tag name
	List() ([]string, error)

	// Obtain return ids from tag name
	Obtain(tag string) (*TagItem, error)
}

// TagItem for tag
type TagItem struct {
	Tag         string `json:"tag"`
	Max         uint64 `json:"max"`
	Min         uint64 `json:"min"`
	Description string `json:"description"`
}

// CommonProvider implement segment provider
type CommonProvider struct {
	mutex sync.RWMutex

	name    string
	storage Storager
	tags    map[string]*TagItem
}

// NewProvider return CommomProvider instance
func NewProvider(name string, storage Storager) (*CommonProvider, error) {
	tags, err := storage.List()
	if err != nil {
		return nil, err
	}

	p := &CommonProvider{
		name:    name,
		storage: storage,
		tags:    make(map[string]*TagItem, len(tags)),
	}
	for _, tag := range tags {
		item, err := storage.Obtain(tag)
		if err != nil {
			return nil, err
		}
		// TODO: check min and max
		p.tags[tag] = item
	}

	return p, nil
}

// Name implement segment.Provider Name
func (p *CommonProvider) Name() string {
	return p.name
}

// Next implement segment.Provider Next
func (p *CommonProvider) Next(arg *segment.Arguments) (*segment.Result, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	item, ok := p.tags[arg.Tag]
	if !ok || item.Min > item.Max {
		var err error
		item, err = p.storage.Obtain(arg.Tag)
		if err != nil {
			return nil, err
		}

		if item.Min > item.Max {
			return nil, fmt.Errorf("min >= max")
		}
		p.tags[arg.Tag] = item
	}

	r := &segment.Result{
		Name: p.name,
		Next: item.Min,
		Labels: map[string]string{
			segment.LabelTag: arg.Tag,
		},
	}
	item.Min++
	return r, nil
}
