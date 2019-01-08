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

package logger

import (
	"fmt"
	"path"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

type callerHook struct {
}

func (callerHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

type source struct {
	p string // package name
	f string // source file name
	m string // function name
	l string // line
}

func (callerHook) Fire(entry *logrus.Entry) error {
	if !isCallerEnable {
		return nil
	}

	c := &source{
		p: "-",
		f: "-",
		m: "-",
		l: "-",
	}

	pc := make([]uintptr, 10, 10)
	cnt := runtime.Callers(7, pc)
	var last = -1
	for i := cnt; i >= 0; i-- {
		fu := runtime.FuncForPC(pc[i-1])
		name := fu.Name()

		if strings.Contains(name, logrusPackageID) {
			break
		}

		last = i
	}

	if last != -1 && last < cnt {
		fu := runtime.FuncForPC(pc[last])
		name := fu.Name()
		file, line := fu.FileLine(pc[last])
		c.p = path.Dir(name)
		c.f = path.Base(file)
		c.m = path.Base(name)
		c.l = fmt.Sprintf("%v", line)
	}

	entry.Data[loggerCallerKeyName] = c
	return nil
}
