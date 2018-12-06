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

package fs

import (
	"encoding/json"
	"sync/atomic"
)

const (
	// SetSuccess defines Set operate success stats
	SetSuccess = iota
	// SetFail defines Set operate fail stats
	SetFail
	// DeleteSuccess defines Delete operate success stats
	DeleteSuccess
	// DeleteFail defines Delete operate fail stats
	DeleteFail
	// CreateSuccess defines Create operate success stats
	CreateSuccess
	// CreateFail defines Create operate fail stats
	CreateFail
	// UpdateSuccess defines Update operate success stats
	UpdateSuccess
	// UpdateFail defines Update operate fail stats
	UpdateFail
	// GetSuccess defines Get operate success stats
	GetSuccess
	// GetFail defines Get operate fail stats
	GetFail
)

// Stater is interface define for stats module
type Stater interface {
	toJSON() []byte
	Inc(int)
	Clone() Stater
}

// defStats struct holds the stats data
type defStats struct {
	GetSuccess uint64 `json:"getSuccess"`
	GetFail    uint64 `json:"getFail"`

	SetSuccess uint64 `json:"setSuccess"`
	SetFail    uint64 `json:"setFail"`

	UpdateSuccess uint64 `json:"updateSuccess"`
	UpdateFail    uint64 `json:"updateFail"`

	CreateSuccess uint64 `json:"createSuccess"`
	CreateFail    uint64 `json:"createFail"`

	DeleteSuccess uint64 `json:"deleteSuccess"`
	DeleteFail    uint64 `json:"deleteFail"`
}

func newStater() Stater {
	return &defStats{}
}

func (s *defStats) Clone() Stater {
	return &defStats{
		GetSuccess:    s.GetSuccess,
		GetFail:       s.GetFail,
		SetSuccess:    s.SetSuccess,
		SetFail:       s.SetFail,
		UpdateSuccess: s.UpdateSuccess,
		UpdateFail:    s.UpdateFail,
		CreateSuccess: s.CreateSuccess,
		CreateFail:    s.CreateFail,
		DeleteSuccess: s.DeleteSuccess,
		DeleteFail:    s.DeleteFail,
	}
}

func (s *defStats) toJSON() []byte {
	b, _ := json.Marshal(s)
	return b
}

// Inc will increment the stats value specified by field define
func (s *defStats) Inc(field int) {
	switch field {
	case SetSuccess:
		atomic.AddUint64(&s.SetSuccess, 1)
	case SetFail:
		atomic.AddUint64(&s.SetFail, 1)
	case GetSuccess:
		atomic.AddUint64(&s.GetSuccess, 1)
	case GetFail:
		atomic.AddUint64(&s.GetFail, 1)
	case UpdateSuccess:
		atomic.AddUint64(&s.UpdateSuccess, 1)
	case UpdateFail:
		atomic.AddUint64(&s.UpdateFail, 1)
	case CreateSuccess:
		atomic.AddUint64(&s.CreateSuccess, 1)
	case CreateFail:
		atomic.AddUint64(&s.CreateFail, 1)
	case DeleteSuccess:
		atomic.AddUint64(&s.DeleteSuccess, 1)
	case DeleteFail:
		atomic.AddUint64(&s.DeleteFail, 1)
	}
}
