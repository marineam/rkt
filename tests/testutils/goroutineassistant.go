// Copyright 2015 The rkt Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package testutils

import (
	"fmt"
	"sync"
	"testing"
)

type GoroutineAssistant struct {
	s  chan string
	wg sync.WaitGroup
	t  *testing.T
}

func NewGoroutineAssistant(t *testing.T) *GoroutineAssistant {
	return &GoroutineAssistant{
		s: make(chan string),
		t: t,
	}
}

func (a *GoroutineAssistant) Fatalf(s string, args ...interface{}) {
	a.s <- fmt.Sprintf(s, args...)
}

func (a *GoroutineAssistant) Add(n int) {
	a.wg.Add(n)
}

func (a *GoroutineAssistant) Done() {
	a.wg.Done()
}

func (a *GoroutineAssistant) Wait() {
	done := make(chan struct{})
	go func() {
		a.wg.Wait()
		done <- struct{}{}
	}()
	select {
	case s := <-a.s:
		a.t.Fatalf(s)
	case <-done:
	}
}
