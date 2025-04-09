package bufreader

import (
	"errors"
	"slices"
	"sync"
)

type Pool interface {
	Alloc(int) []byte
	Free([]byte)
}

var _ Pool = (*SyncPool)(nil)

var (
	// ErrInvalidSize when the input size parameters are invalid
	ErrInvalidSize     = errors.New("invalid size parameters")
	ErrMinSizeTooSmall = errors.New("minimum size must be positive")
	ErrMaxSizeTooSmall = errors.New("maximum size must be greater than or equal to minimum size")
	ErrFactorTooSmall  = errors.New("growth factor must be positive")
)

// SyncPool is a sync.Pool based slab allocation memory pool
type SyncPool struct {
	classes     []sync.Pool // different size memory pools
	classesSize []int       // size of each memory pool
	minSize     int         // minimum chunk size
	maxSize     int         // maximum chunk size
	sizeLookup  []uint32    // fast lookup table from size to class index
}

// NewSyncPool creates a sync.Pool based slab allocation memory pool
// minSize: minimum chunk size
// maxSize: maximum chunk size
// factor: factor for controlling chunk size growth
func NewSyncPool(minSize, maxSize, factor int) (*SyncPool, error) {
	if minSize <= 0 {
		return nil, ErrMinSizeTooSmall
	}
	if maxSize < minSize {
		return nil, ErrMaxSizeTooSmall
	}
	if factor <= 0 {
		return nil, ErrFactorTooSmall
	}

	var classesSize []int
	var curSize = minSize

	for curSize < maxSize {
		classesSize = append(classesSize, curSize)
		nextSize := int(float64(curSize) * (1.0 + 1.0/float64(factor)))
		if nextSize <= curSize {
			nextSize = curSize + minSize // make sure at least grow by minSize
		}
		curSize = nextSize
	}
	classesSize = append(classesSize, maxSize)

	n := len(classesSize)

	pool := &SyncPool{
		classes:     make([]sync.Pool, n),
		classesSize: classesSize,
		minSize:     minSize,
		maxSize:     maxSize,
		sizeLookup:  make([]uint32, maxSize+1),
	}

	// initialize each memory pool class and fill the lookup table
	for k := range pool.classes {
		size := pool.classesSize[k]
		pool.classes[k].New = func() interface{} {
			buf := make([]byte, size)
			return &buf
		}

		// fill the lookup table
		start := 0
		if k > 0 {
			start = pool.classesSize[k-1] + 1
		}
		for i := start; i <= size && i <= maxSize; i++ {
			pool.sizeLookup[i] = uint32(k)
		}
	}

	return pool, nil
}

// Alloc allocates a []byte from the internal slab class
// if there is no free block, it will create a new one
func (pool *SyncPool) Alloc(size int) []byte {
	if size <= 0 {
		return make([]byte, 0)
	}

	if size <= pool.maxSize {
		classIndex := pool.sizeLookup[size]
		mem := pool.classes[classIndex].Get().(*[]byte)
		return (*mem)[:size]
	}

	return make([]byte, size)
}

// Free frees the []byte allocated from Pool.Alloc
func (pool *SyncPool) Free(mem []byte) {
	if len(mem) == 0 {
		return
	}

	size := cap(mem)
	// for memory blocks less than pool.minSize or larger than pool.maxSize, let GC handle it
	if size < pool.minSize || size > pool.maxSize {
		return
	}

	classIndex := pool.sizeLookup[size]
	// reset the slice to avoid memory leaks
	mem = mem[:cap(mem)]
	// zero out sensitive data
	for i := range mem {
		mem[i] = 0
	}
	pool.classes[classIndex].Put(&mem)
}

// ClassCount returns the number of memory pool classes
func (pool *SyncPool) ClassCount() int {
	return len(pool.classes)
}

// ClassSizes returns the sizes of all memory pool classes
func (pool *SyncPool) ClassSizes() []int {
	return slices.Clone(pool.classesSize)
}

// MinMaxSize returns the minimum and maximum block sizes of the pool
func (pool *SyncPool) MinMaxSize() (min, max int) {
	return pool.minSize, pool.maxSize
}
