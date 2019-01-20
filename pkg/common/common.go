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
	"fmt"
	"strings"
)

// ProviderResult for provider next
type ProviderResult struct {
	Name   string            `json:"name"`
	Next   uint64            `json:"next"`
	Labels map[string]string `json:"labels,omitempty"`
}

// Version for fyllo version profile
type Version struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Commit      string `json:"commit"`
	Description string `json:"description"`
}

// SplitKeyValueArrayString split the string to map, the string format should be:
// key1=value1,key2=value2
// the different key is split by comma, and the key/value is split by equal sign
// if there is two same key name, it will return an error
func SplitKeyValueArrayString(args string) (map[string]string, error) {
	return SplitKeyValueArrayStringSep(args, ",")
}

// SplitKeyValueArrayStringSep split the string to map, the string format should be:
// key1=value1[sep]key2=value2
// the different key is split by sep, and the key/value is split by equal sign
// if there is two same key name, it will return an error
func SplitKeyValueArrayStringSep(args string, sep string) (map[string]string, error) {
	kvs := strings.Split(args, sep)
	r := make(map[string]string)
	for _, field := range kvs {
		kv := strings.Split(field, "=")
		if 2 != len(kv) {
			return nil, fmt.Errorf("the field[%s] doesn't match key=value format", field)
		}

		key := kv[0]
		if prev, ok := r[key]; ok {
			return nil, fmt.Errorf("the field[%s] key[%s] is duplicated, prev value[%s]", field, key, prev)
		}
		r[key] = kv[1]
	}
	return r, nil
}
