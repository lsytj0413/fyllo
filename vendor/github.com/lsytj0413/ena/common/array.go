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
	"errors"
	"reflect"
)

// ReverseSlice will reverse a array or slice
func ReverseSlice(result interface{}) error {
	resultv := reflect.ValueOf(result)
	if resultv.Kind() != reflect.Ptr || (resultv.Elem().Kind() != reflect.Slice && resultv.Elem().Kind() != reflect.Array) {
		return errors.New("result argument must be a slice address or a array address")
	}

	slicev := resultv.Elem()
	for left, right := 0, slicev.Len()-1; left < right; left, right = left+1, right-1 {
		leftV, rightV := slicev.Index(left), slicev.Index(right)
		tmp := leftV.Interface()

		leftV.Set(rightV)
		rightV.Set(reflect.ValueOf(tmp))
	}

	return nil
}
