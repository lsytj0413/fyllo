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
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/lsytj0413/fyllo/pkg/errors"
	"github.com/lsytj0413/fyllo/pkg/segment"
)

type commonProviderTestSuite struct {
	suite.Suite
}

func (s *commonProviderTestSuite) TestNewProviderOk() {
	tagItems := []*TagItem{
		{
			Tag:         "1",
			Max:         10,
			Min:         1,
			Description: "1",
		},
		{
			Tag:         "2",
			Max:         20,
			Min:         2,
			Description: "2",
		},
	}
	mockStorager := newMockStorager(tagItems)

	providerName := "test"
	p, err := NewProvider(providerName, mockStorager)
	s.NoError(err)
	s.NotNil(p)

	s.Equal(p.Name(), providerName)
	s.Equal(len(tagItems), len(p.tags))

	for _, item := range tagItems {
		v := p.tags[item.Tag]
		if !reflect.DeepEqual(item, v) {
			s.FailNowf("item equal failed", "expect %v, got %v", item, v)
		}
	}
}

func (s *commonProviderTestSuite) TestNewProviderListFailed() {
	errString := "List Failed"
	mockStorager := &mockStorageObject{}
	mockStorager.On("List").Return([]string{}, fmt.Errorf(errString))

	providerName := "test"
	p, err := NewProvider(providerName, mockStorager)
	s.Nil(p)
	if err.Error() != errString {
		s.FailNowf("error string failed", "expect %v, got %v", errString, err)
	}
}

func (s *commonProviderTestSuite) TestNewProviderObtainFailed() {
	tagItems := []*TagItem{
		{
			Tag:         "1",
			Max:         10,
			Min:         1,
			Description: "1",
		},
	}
	mockStorager := newMockStorager(tagItems)
	errString := "Obtain Failed"
	for _, item := range tagItems {
		var nilItem *TagItem
		mockStorager.On("Obtain", item.Tag).Return(nilItem, fmt.Errorf(errString))
	}
	mockStorager.ReverseCall()

	providerName := "test"
	p, err := NewProvider(providerName, mockStorager)
	s.Nil(p)
	if err.Error() != errString {
		s.FailNowf("error string failed", "expect %v, got %v", errString, err)
	}
}

func (s *commonProviderTestSuite) TestNextOk() {
	tagItems := []*TagItem{
		{
			Tag:         "1",
			Max:         10,
			Min:         1,
			Description: "1",
		},
	}
	mockStorager := newMockStorager(tagItems)

	providerName := "test"
	p, err := NewProvider(providerName, mockStorager)
	s.NoError(err)
	s.NotNil(p)

	expect, tagName := uint64(1), tagItems[0].Tag
	r, err := p.Next(&segment.Arguments{
		Tag: tagName,
	})
	s.NoError(err)
	s.Equal(expect, r.Next)
	s.Equal(providerName, r.Name)
	s.Equal(tagName, r.Labels[segment.LabelTag])
}

func (s *commonProviderTestSuite) TestNextMultiValueOk() {
	tagItems := []*TagItem{
		{
			Tag:         "1",
			Max:         5,
			Min:         1,
			Description: "1",
		},
		{
			Tag:         "2",
			Max:         4,
			Min:         2,
			Description: "2",
		},
	}
	mockStorager := newMockStorager(tagItems)

	providerName := "test"
	p, err := NewProvider(providerName, mockStorager)
	s.NoError(err)
	s.NotNil(p)

	s.Equal(p.Name(), providerName)

	type testCase struct {
		description string
		tag         string
		expects     []uint64
	}
	testCases := []testCase{
		{
			description: "normal test tag 1",
			tag:         "1",
			expects:     []uint64{1, 2, 3, 4, 5, 1, 2, 3, 4, 5, 1},
		},
		{
			description: "normal test tag 2",
			tag:         "2",
			expects:     []uint64{2, 3, 4, 2, 3, 4, 2},
		},
	}
	for _, tc := range testCases {
		for _, expect := range tc.expects {
			r, err := p.Next(&segment.Arguments{
				Tag: tc.tag,
			})
			s.NoError(err)
			s.Equal(expect, r.Next, tc.description)
		}
	}
}

func (s *commonProviderTestSuite) TestNextArgNilFailed() {
	tagItems := []*TagItem{
		{
			Tag:         "1",
			Max:         5,
			Min:         1,
			Description: "1",
		},
		{
			Tag:         "2",
			Max:         4,
			Min:         2,
			Description: "2",
		},
	}
	mockStorager := newMockStorager(tagItems)

	providerName := "test"
	p, err := NewProvider(providerName, mockStorager)
	s.NoError(err)
	s.NotNil(p)

	s.Equal(p.Name(), providerName)

	r, err := p.Next(nil)
	s.Nil(r)
	s.Error(err)
}

func (s *commonProviderTestSuite) TestNextObtainFailed() {
	tagItems := []*TagItem{
		{
			Tag:         "1",
			Max:         5,
			Min:         1,
			Description: "1",
		},
		{
			Tag:         "2",
			Max:         4,
			Min:         2,
			Description: "2",
		},
	}
	mockStorager := newMockStorager(tagItems)
	providerName := "test"
	p, err := NewProvider(providerName, mockStorager)
	s.NoError(err)
	s.NotNil(p)

	errString := "Obtain Failed"
	for _, item := range tagItems {
		var nilItem *TagItem
		mockStorager.On("Obtain", item.Tag).Return(nilItem, fmt.Errorf(errString))
	}
	mockStorager.ReverseCall()

	s.Equal(p.Name(), providerName)

	r, err := p.Next(nil)
	s.Nil(r)
	s.Error(err)
}

func (s *commonProviderTestSuite) TestObtainItemOk() {
	tagItems := []*TagItem{
		{
			Tag:         "1",
			Max:         10,
			Min:         1,
			Description: "1",
		},
	}
	mockStorager := newMockStorager(tagItems)

	for _, item := range tagItems {
		r, err := obtainTagNextItem(mockStorager, item.Tag)
		s.NoError(err)

		if !reflect.DeepEqual(item, r) {
			s.FailNowf("item equal failed", "expect %v, got %v", item, r)
		}
		if item == r {
			s.FailNowf("item pointer failed", "should not equal, expect[%p], got[%p]", item, r)
		}
	}
}

func (s *commonProviderTestSuite) TestObtainItemFailed() {
	tagItems := []*TagItem{
		{
			Tag:         "1",
			Max:         10,
			Min:         11,
			Description: "1",
		},
	}
	mockStorager := newMockStorager(tagItems)
	errString := "Obtain Failed"
	for _, item := range tagItems {
		var tagItem *TagItem
		mockStorager.On("Obtain", item.Tag).Return(tagItem, fmt.Errorf(errString))
	}
	mockStorager.ReverseCall()

	for _, item := range tagItems {
		r, err := obtainTagNextItem(mockStorager, item.Tag)
		s.Error(err)
		s.Nil(r)
	}
}

func (s *commonProviderTestSuite) TestObtainItemValueFailed() {
	tagItems := []*TagItem{
		{
			Tag:         "1",
			Max:         10,
			Min:         11,
			Description: "1",
		},
	}
	mockStorager := newMockStorager(tagItems)

	for _, item := range tagItems {
		r, err := obtainTagNextItem(mockStorager, item.Tag)
		s.Error(err)
		if !errors.Is(err, errors.EcodeSegmentRangeFailed) {
			s.Failf("error code failed", "expect[%v], got[%v]", errors.EcodeSegmentRangeFailed, err)
		}
		s.Nil(r)
	}
}

func TestCommonProviderTestSuite(t *testing.T) {
	s := &commonProviderTestSuite{}
	suite.Run(t, s)
}
