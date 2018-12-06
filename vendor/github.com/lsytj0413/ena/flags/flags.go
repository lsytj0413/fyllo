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

package flags

import (
	"flag"
	"fmt"
	"io"
)

// DeprecatedFlag encapsulates a flag that may have been previously valid but
// is now deprecated. If a DeprecatedFlag is set, an error occurs.
type DeprecatedFlag struct {
	Name string
}

// Set flag value
func (f *DeprecatedFlag) Set(_ string) error {
	return fmt.Errorf(`flag "%s" is no longer supported`, f.Name)
}

func (f *DeprecatedFlag) String() string {
	return ""
}

// IgnoredFlag encapsulates a flag that may have been previously valid but is
// now ignored. If an IgnoredFlag is set, a warning is printed and operation continues.
type IgnoredFlag struct {
	Name string
	w    io.Writer
}

// IsBoolFlag is defined to allow the flag to be defined without an argument
func (f *IgnoredFlag) IsBoolFlag() bool {
	return true
}

// Set flag value
func (f *IgnoredFlag) Set(_ string) error {
	fmt.Fprintf(f.w, `flag "%s" is no longer supported - ignoring`, f.Name)
	return nil
}

func (f *IgnoredFlag) String() string {
	return ""
}

// IsSet check the provide flag is set by caller
func IsSet(fs *flag.FlagSet, name string) bool {
	set := false
	fs.Visit(func(f *flag.Flag) {
		if f.Name == name {
			set = true
		}
	})

	return set
}
