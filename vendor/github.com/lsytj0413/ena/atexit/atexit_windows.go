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

// +build windows

// Package atexit provide a way to run Handler at process exit
package atexit

import "os"

// Handler is the function type which will be called at exit
type Handler func()

// RegisterHandler add handler which will be called at exit
func RegisterHandler(handler Handler) {}

// HandleInterrupts is a no-op on windows
func HandleInterrupts() {}

// Exit calls os.Exit
func Exit(code int) {
	os.Exit(code)
}
