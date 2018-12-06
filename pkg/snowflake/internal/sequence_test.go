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

package internal

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/lsytj0413/ena/cerror"
	ierror "github.com/lsytj0413/fyllo/pkg/error"
)

type sequenceTestSuite struct {
	suite.Suite
}

func (s *sequenceTestSuite) TestNewSequenceOk() {
	maxes := []uint64{100, 20, 30, 40}
	for _, max := range maxes {
		sq := newSequence(max).(*defSequence)
		s.Equal(sq.now, uint64(0))
		s.Equal(sq.max, max)
	}
}

func (s *sequenceTestSuite) TestResetOk() {
	maxes := []uint64{100, 20, 30, 40}
	for _, max := range maxes {
		sq := newSequence(max).(*defSequence)
		sq.now = max
		sq.reset()
		s.Equal(sq.now, uint64(0))
	}
}

func (s *sequenceTestSuite) TestIsOutRange() {
	maxes := []uint64{100, 20, 30, 40}
	for _, max := range maxes {
		sq := newSequence(max).(*defSequence)

		for i := 0; max < max; i++ {
			sq.now = uint64(i)
			s.Equal(false, sq.isOutRange())
		}

		for i := max; i < max+10; i++ {
			sq.now = uint64(i)
			s.Equal(true, sq.isOutRange())
		}
	}
}

func (s *sequenceTestSuite) TestNextOk() {
	max := uint64(100)
	sq := newSequence(max).(*defSequence)
	for i := uint64(0); i < max; i++ {
		n, err := sq.next()
		s.NoError(err)
		s.Equal(i, n)
	}

	// should return error
	n, err := sq.next()
	s.Error(err)
	s.Equal(uint64(0), n)

	cerr := err.(*cerror.Error)
	s.Equal(ierror.EcodeSequenceOutOfRange, cerr.ErrorCode)
}

func TestSequenceTestSuite(t *testing.T) {
	p := &sequenceTestSuite{}
	suite.Run(t, p)
}
