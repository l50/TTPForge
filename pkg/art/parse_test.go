/*
Copyright © 2023-present, Meta Platforms, Inc. and affiliates
Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:
The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.
THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package art_test

import (
	"testing"

	"github.com/facebookincubator/ttpforge/pkg/art"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name   string
		path   string
		expect *art.Config
		err    error
	}{
		// TODO: Add test cases here
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// cfg, err := art.NewConfig(tc.path)
			// if err != tc.err {
			// 	t.Errorf("Expected error %v, but got %v", tc.err, err)
			// }
			// Add more assertions
		})
	}
}

func TestLoadArtYAML(t *testing.T) {
	tests := []struct {
		name   string
		path   string
		expect *art.Atomic
		err    error
	}{
		// TODO: Add test cases here
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			atomic := art.NewAtomic()
			err := atomic.LoadArtYAML(tc.path)
			if err != tc.err {
				t.Errorf("Expected error %v, but got %v", tc.err, err)
			}
			// Add more assertions
		})
	}
}

func TestNewAbility(t *testing.T) {
	tests := []struct {
		name     string
		id       int64
		command  string
		expected *art.Ability
	}{
		// TODO: Add test cases here
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// actual := art.NewAbility(tc.id, tc.command)
			// Assertions, e.g., compare actual and expected Ability
		})
	}
}

func TestNewVar(t *testing.T) {
	tests := []struct {
		name     string
		id       int64
		nameVar  string
		value    string
		expected *art.Var
	}{
		// TODO: Add test cases here
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// actual := art.NewVar(tc.id, tc.nameVar, tc.value)
			// Assertions, e.g., compare actual and expected Var
		})
	}
}
