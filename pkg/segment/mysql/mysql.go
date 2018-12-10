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

package mysql

import (
	"github.com/lsytj0413/fyllo/pkg/segment"
	"github.com/lsytj0413/fyllo/pkg/segment/internal"
)

const (
	// ProviderName for the mysql segment provider
	ProviderName = "mysql"
)

// Options is mysql segment provider option
type Options struct {
	Args string
}

type mysqlStorage struct {
}

func (m *mysqlStorage) List() ([]string, error) {
	return []string{}, nil
}

func (m *mysqlStorage) Obtain(tag string) (*internal.TagItem, error) {
	return nil, nil
}

// NewProvider return mysql segment provider implement
func NewProvider(options *Options) (segment.Provider, error) {
	return internal.NewProvider(ProviderName, &mysqlStorage{})
}
