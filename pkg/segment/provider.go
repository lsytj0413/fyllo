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

import "github.com/lsytj0413/fyllo/pkg/common"

// Provider for provide segment id
type Provider interface {
	Name() string
	Next(arg *Arguments) (*Result, error)
}

// Result for segment Next
type Result = common.ProviderResult

// Arguments for segment generate
type Arguments struct {
	Tag string `json:"tag"`
}

const (
	// LabelTag for tag label key
	LabelTag = "tag"
)
