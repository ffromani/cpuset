/*
Copyright 2020 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cpuset

import (
	"reflect"
	"testing"
)

func TestEmpty(t *testing.T) {
	checkEmpty(t, Empty())
}

func TestParseEmpty(t *testing.T) {
	cpus, err := Parse("")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	checkEmpty(t, cpus)
}

func TestParseRange(t *testing.T) {
	testCases := []tcase{
		{"0,1,2,3", []int{0, 1, 2, 3}},
		// we don't elide dupes
		{"0,1,1,3", []int{0, 1, 1, 3}},
	}
	for _, testCase := range testCases {
		testCase.CheckSlices(t)
	}
}

func TestRangeInterval(t *testing.T) {
	testCases := []tcase{
		{"0-7", []int{0, 1, 2, 3, 4, 5, 6, 7}},
		{"0-3,4-7", []int{0, 1, 2, 3, 4, 5, 6, 7}},
		{"1,2-5,6", []int{1, 2, 3, 4, 5, 6}},
	}
	for _, testCase := range testCases {
		testCase.CheckSlices(t)
	}
}

func TestRangeIntervalMalformed(t *testing.T) {
	testCases := []tcase{
		{",", nil},
		{"-", nil},
		{"-,", nil},
		{",-", nil},
		{",-,", nil},
		{",,", nil},
		{"1,2-,6", nil},
		{"1,-3,6", nil},
		{"1,2-4-6,8", nil},
		{"1,-,8", nil},
	}
	for _, testCase := range testCases {
		testCase.ExpectError(t)
	}
}

func checkEmpty(t *testing.T, cpus []int) {
	if cpus == nil {
		t.Errorf("empty must not be nil")
	}
	if len(cpus) != 0 {
		t.Errorf("empty must have zero length")
	}

}

type tcase struct {
	s string
	v []int
}

func (tc tcase) CheckSlices(t *testing.T) {
	cpus, err := Parse(tc.s)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(cpus, tc.v) {
		t.Errorf("slices differ: expected=%v cpus=%v", tc.v, cpus)
	}
}

func (tc tcase) ExpectError(t *testing.T) {
	_, err := Parse(tc.s)
	if err == nil {
		t.Errorf("unexpectedly ok: %q", tc.s)
	}
}
