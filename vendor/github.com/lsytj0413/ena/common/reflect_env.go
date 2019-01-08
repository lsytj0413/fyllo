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
	"os"
	"reflect"
	"regexp"
)

var (
	getenv   func(string) (string, bool)
	submatch func(string) (string, bool)

	reg *regexp.Regexp
)

func replaceEnvImpl(v reflect.Value) error {
	switch v.Kind() {
	case reflect.String:
		replace(v)
	case reflect.Struct:
		n := v.NumField()
		for i := 0; i < n; i++ {
			field := v.Field(i)
			if field.Kind() == reflect.String || field.Kind() == reflect.Struct {
				if err := replaceEnvImpl(field); err != nil {
					return err
				}
			}
		}
	default:
		return errors.New("value argument must be a struct/string address")
	}

	return nil
}

func replace(v reflect.Value) {
	env, exists := submatch(v.String())
	if exists {
		env, exists = getenv(env)
		if exists {
			v.SetString(env)
		}
	}
}

// ReplaceEnv will fill the value element with enviroment value
func ReplaceEnv(value interface{}) error {
	resultv := reflect.ValueOf(value)
	if resultv.Kind() != reflect.Ptr {
		return errors.New("value argument must be a ptr")
	}

	return replaceEnvImpl(resultv.Elem())
}

func init() {
	getenv = func(key string) (v string, exists bool) {
		v = os.Getenv(key)
		if v != "" {
			exists = true
		}
		return
	}

	reg = regexp.MustCompile(`^\${(.+)}$`)
	submatch = func(key string) (v string, exists bool) {
		values := reg.FindStringSubmatch(key)
		if len(values) == 2 {
			v, exists = values[1], true
		}
		return
	}
}
