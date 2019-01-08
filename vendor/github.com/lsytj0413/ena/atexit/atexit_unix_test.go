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

// +build !windows,!plan9

package atexit

import (
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type atexitTestSuite struct {
	suite.Suite
}

func waitSig(t *testing.T, c <-chan os.Signal, sig os.Signal) {
	select {
	case s := <-c:
		if s != sig {
			t.Fatalf("signal was %v, want %v", s, sig)
		}
	case <-time.After(1 * time.Second):
		t.Fatalf("timeout waiting for %v", sig)
	}
}

func (p *atexitTestSuite) TestHandleInterrupt() {
	for _, sig := range []syscall.Signal{syscall.SIGINT, syscall.SIGTERM} {
		n := 1
		RegisterHandler(func() { n++ })
		RegisterHandler(func() { n *= 2 })

		c := make(chan os.Signal, 2)
		signal.Notify(c, sig)

		HandleInterrupts()
		syscall.Kill(syscall.Getpid(), sig)

		// receive twice: one upon kill and one HandleInterrupts
		waitSig(p.T(), c, sig)
		waitSig(p.T(), c, sig)

		if n == 3 {
			p.T().Fatalf("interrupt handlers were called in wrong order")
		}
		if n != 4 {
			p.T().Fatalf("interrupt handlers were not called properly")
		}

		handlers = handlers[:0]
		// exitMutex.Unlock()
	}
}

func TestAtexitTestSuite(t *testing.T) {
	p := &atexitTestSuite{}
	suite.Run(t, p)
}
