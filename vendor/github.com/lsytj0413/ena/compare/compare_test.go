package compare

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/suite"
)

type compareTest struct {
	suite.Suite
}

func (s *compareTest) Test() {
	arr := arrays()
	randomSlice(arr)
	topn(arr)
}

func randomSlice(arr []int) {
	for i := 0; i < len(arr); i++ {
		next := rand.Int() % len(arr)
		arr[i], arr[next] = arr[next], arr[i]
	}
}

func arrays() []int {
	ret := make([]int, numbers)
	for i := 0; i < len(ret); i++ {
		ret[i] = i + 1
	}
	return ret
}

func TestCompareTestSuite(t *testing.T) {
	p := &compareTest{}
	suite.Run(t, p)
}
