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
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

type commonTestSuite struct {
	suite.Suite
}

func (s *commonTestSuite) TestSplitKeyValueArrayStringSepOk() {
	type testCase struct {
		description string
		args        string
		sep         string
		target      map[string]string
	}
	testCases := []testCase{
		{
			description: "01: one kv",
			args:        "key=value",
			sep:         ";",
			target: map[string]string{
				"key": "value",
			},
		},
		{
			description: "02: multi kv",
			args:        "key=value;key1=value1;key2=value2",
			sep:         ";",
			target: map[string]string{
				"key":  "value",
				"key1": "value1",
				"key2": "value2",
			},
		},
		{
			description: "02: comma sep",
			args:        "key=value,key1=value1,key2=value2",
			sep:         ",",
			target: map[string]string{
				"key":  "value",
				"key1": "value1",
				"key2": "value2",
			},
		},
	}
	for _, tc := range testCases {
		actual, err := SplitKeyValueArrayStringSep(tc.args, tc.sep)
		if err != nil {
			s.FailNowf(tc.description, "expect no error but failed, %v", err)
		}
		s.Equal(tc.target, actual, tc.description)
	}
}

func (s *commonTestSuite) TestSplitKeyValueArrayStringSepFailed() {
	matcherr, duperr := "doesn't match", "is duplicated"
	type testCase struct {
		description string
		args        string
		sep         string
		err         string
	}
	testCases := []testCase{
		{
			description: "01: kv match failed",
			args:        "keyvalue",
			sep:         ";",
			err:         matcherr,
		},
		{
			description: "02: duplicate key",
			args:        "key=value;key=value1;key2=value2",
			sep:         ";",
			err:         duperr,
		},
	}
	for _, tc := range testCases {
		actual, err := SplitKeyValueArrayStringSep(tc.args, tc.sep)
		if err == nil {
			s.FailNowf(tc.description, "expect error but got nil")
		}
		s.Nil(actual)
		if !strings.Contains(err.Error(), tc.err) {
			s.FailNowf(tc.description, "expect %s err got %v", tc.err, err)
		}
	}
}

func TestCommonTestSuite(t *testing.T) {
	s := &commonTestSuite{}
	suite.Run(t, s)
}
