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

// Package atexit provide a way to run Handler at process exit
package atexit

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// Handler is the function type which will be called at exit
type Handler func()

var (
	handlers      = []Handler{}
	registerMutex sync.Mutex
	exitMutex     sync.Mutex

	ctx, cancel = context.WithCancel(context.Background())

	// ErrWriter for error log, os.Stderr is default
	ErrWriter io.Writer
	// StdWriter for normal log, os.Stdout is default
	StdWriter io.Writer
)

func runHandler(handler Handler) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintf(ErrWriter, "atexit: error on hander, %v", err)
		}
	}()

	handler()
}

func runHandlers(h []Handler) {
	for _, handler := range h {
		runHandler(handler)
	}
}

// Exit run all handler and terminates the program by os.Exit
func Exit(code int) {
	exitMutex.Lock()
	defer exitMutex.Unlock()

	cancel()
	runHandlers(handlers)

	os.Exit(code)
}

// RegisterHandler add handler which will be called at exit
func RegisterHandler(handler Handler) {
	registerMutex.Lock()
	defer registerMutex.Unlock()

	handlers = append(handlers, handler)
}

// HandleInterrupts will calls the handler functions on receiving a SIGINT or SIGTERM
func HandleInterrupts() {
	notifier := make(chan os.Signal, 1)
	signal.Notify(notifier, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		var sig os.Signal
		select {
		case <-ctx.Done():
			return
		case sig = <-notifier:
		}

		tmpHandlers := func() []Handler {
			registerMutex.Lock()
			defer registerMutex.Unlock()

			tmp := make([]Handler, len(handlers))
			copy(tmp, handlers)
			return tmp
		}()

		exitMutex.Lock()
		defer exitMutex.Unlock()

		fmt.Fprintf(StdWriter, "atexit: received %v signal, shutting down...", sig)

		runHandlers(tmpHandlers)

		signal.Stop(notifier)
		pid := syscall.Getpid()

		// exit directly if it is the "init" process, since the kernel will not help to kill pid 1.
		if pid == 1 {
			os.Exit(0)
		}

		err := syscall.Kill(pid, sig.(syscall.Signal))
		if err != nil {
			fmt.Fprintf(ErrWriter, "atexit: kill %v error, %v", pid, err)
		}
	}()
}

func init() {
	ErrWriter = os.Stderr
	StdWriter = os.Stdout
}
