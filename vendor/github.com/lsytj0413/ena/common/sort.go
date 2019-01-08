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

package common

import (
	"sort"
)

// Lessor is a function witch compare data's i and j element
type Lessor = func(i int, j int) bool

// Swaper is a function witch swap data's i and j element
type Swaper = func(i int, j int)

// Lener is a function witch returns data's length
type Lener = func() int

type sortHelper struct {
	lessor Lessor
	swaper Swaper
	lener  Lener
}

func (s sortHelper) Less(i int, j int) bool {
	return s.lessor(i, j)
}

func (s sortHelper) Swap(i int, j int) {
	s.swaper(i, j)
}

func (s sortHelper) Len() int {
	return s.lener()
}

// Sort is a helper function for sort.Sort
func Sort(lessor Lessor, lener Lener, swaper Swaper) {
	sorter := &sortHelper{
		lessor: lessor,
		swaper: swaper,
		lener:  lener,
	}
	sort.Sort(sorter)
}
