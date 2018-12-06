package conc

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// ConcurrentArray is thread-safe array
type ConcurrentArray interface {
	Set(index uint32, elem interface{}) (err error)
	Get(index uint32) (elem interface{}, err error)
	Len() uint32
}

const (
	maxArrayLength = 100
)

type defaultArray struct {
	mutex  sync.Mutex
	length uint32
	val    atomic.Value
}

// NewConcurrentArray will create a ConcurrentArray instance
func NewConcurrentArray(length uint32) (ConcurrentArray, error) {
	if length > maxArrayLength {
		return nil, fmt.Errorf("NewConcurrentArray: length must be less than [%v]", maxArrayLength+1)
	}

	array := defaultArray{}
	array.length = length
	array.val.Store(make([]interface{}, array.length))
	return &array, nil
}

func (array *defaultArray) Set(index uint32, elem interface{}) (err error) {
	if err = array.checkIndex(index); err != nil {
		return
	}

	array.mutex.Lock()
	defer array.mutex.Unlock()

	newArray := make([]interface{}, array.length)
	copy(newArray, array.val.Load().([]interface{}))
	newArray[index] = elem
	array.val.Store(newArray)
	return
}

func (array *defaultArray) Get(index uint32) (elem interface{}, err error) {
	if err = array.checkIndex(index); err != nil {
		return
	}

	elem = array.val.Load().([]interface{})[index]
	return
}

func (array *defaultArray) Len() uint32 {
	return array.length
}

func (array *defaultArray) checkIndex(index uint32) error {
	if index >= array.length {
		return fmt.Errorf("ConcurrentArray.checkIndex: Index out of range [0, %d)", array.length)
	}
	return nil
}
