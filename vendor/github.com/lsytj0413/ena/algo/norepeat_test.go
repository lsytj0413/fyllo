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

package algo

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type norepeatTestSuite struct {
	suite.Suite
}

type norepeats struct {
	value string
	no    string
	len   int
}

func (s *norepeatTestSuite) TestEmpty() {
	v := norepeats{
		value: "",
		no:    "",
		len:   0,
	}

	v0, l := FindNoRepeat(v.value)
	s.Equal(v.no, v0)
	s.Equal(v.len, l)
}

func (s *norepeatTestSuite) TestLen1() {
	vs := []norepeats{
		{
			value: "a",
			no:    "a",
			len:   1,
		},
		{
			value: "z",
			no:    "z",
			len:   1,
		},
		{
			value: "Z",
			no:    "Z",
			len:   1,
		},
		{
			value: "8",
			no:    "8",
			len:   1,
		},
		{
			value: "中",
			no:    "中",
			len:   1,
		},
	}

	for _, v := range vs {
		v0, l := FindNoRepeat(v.value)
		s.Equal(v.no, v0)
		s.Equal(v.len, l)
	}
}

func (s *norepeatTestSuite) TestLen2() {
	vs := []norepeats{
		{
			value: "aa",
			no:    "a",
			len:   1,
		},
		{
			value: "za",
			no:    "za",
			len:   2,
		},
		{
			value: "ZZ",
			no:    "Z",
			len:   1,
		},
		{
			value: "87",
			no:    "87",
			len:   2,
		},
		{
			value: "中国",
			no:    "中国",
			len:   2,
		},
		{
			value: "中中",
			no:    "中",
			len:   1,
		},
		{
			value: "a中",
			no:    "a中",
			len:   2,
		},
	}

	for _, v := range vs {
		v0, l := FindNoRepeat(v.value)
		s.Equal(v.no, v0)
		s.Equal(v.len, l)
	}
}

func (s *norepeatTestSuite) TestLenMulti() {
	vs := []norepeats{
		{
			value: "aba",
			no:    "ab",
			len:   2,
		},
		{
			value: "zasfds",
			no:    "zasfd",
			len:   5,
		},
		{
			value: "ABDEFGACBEF",
			no:    "BDEFGAC",
			len:   7,
		},
		{
			value: "abcdefghij",
			no:    "abcdefghij",
			len:   10,
		},
		{
			value: "中国三ss复",
			no:    "中国三s",
			len:   4,
		},
		{
			value: "wernsdflkjijh3434",
			no:    "wernsdflkji",
			len:   11,
		},
		{
			value: "aabcceddabc",
			no:    "dabc",
			len:   4,
		},
		{
			value: "abcabcbb",
			no:    "abc",
			len:   3,
		},
		{
			value: "bbbbb",
			no:    "b",
			len:   1,
		},
		{
			value: "pwwkew",
			no:    "wke",
			len:   3,
		},
		{
			value: "qpxrjxkltzyx",
			no:    "rjxkltzy",
			len:   8,
		},
	}

	for _, v := range vs {
		v0, l := FindNoRepeat(v.value)
		s.Equal(v.no, v0)
		s.Equal(v.len, l)
	}
}

func TestNoRepeatTestSuite(t *testing.T) {
	s := &norepeatTestSuite{}
	suite.Run(t, s)
}
