package slab

import (
	"unsafe"

	"github.com/pkg/errors"
)

// ChanPool is a chan based slab allocation memory pool.
type ChanPool struct {
	classes []chanClass
	minSize int
	maxSize int
	errChan chan error
}

// newChanPool create a chan based slab allocation memory pool.
// minSize is the smallest chunk size.
// maxSize is the lagest chunk size.
// factor is used to control growth of chunk size.
// pageSize is the memory size of each slab class.
func newChanPool(minSize, maxSize, factor int, pageSize int) *ChanPool {
	var (
		chunkSize int = minSize
		n         int = 0
		//pageSize  int = 8192 // 8kb
	)
	for ; chunkSize <= maxSize && chunkSize <= pageSize; chunkSize *= factor {
		n++
	}
	if maxSize > int(chunkSize/factor) && maxSize <= pageSize {
		n++
	}
	pool := &ChanPool{
		classes: make([]chanClass, n),
		minSize: minSize,
		maxSize: maxSize,
	}
	chunkSize = minSize
	for i := 0; i < n; i++ {
		pool.classes[i].size = chunkSize
		pool.classes[i].page = make([]byte, pageSize)
		pool.classes[i].chunks = make(chan []byte, pageSize/chunkSize)
		pool.classes[i].chanPool = pool
		pool.classes[i].pageBegin = uintptr(unsafe.Pointer(&pool.classes[i].page[0]))
		for j := 0; j < pageSize/chunkSize; j++ {
			// lock down the capacity to protect append operation
			mem := pool.classes[i].page[j*chunkSize : (j+1)*chunkSize : (j+1)*chunkSize]
			pool.classes[i].chunks <- mem
			if j == len(pool.classes[i].chunks)-1 {
				pool.classes[i].pageEnd = uintptr(unsafe.Pointer(&mem[0]))
			}
		}

		chunkSize *= factor
	}
	pool.errChan = make(chan error, n*5)
	return pool
}

func (pool *ChanPool) ErrChan() <-chan error {
	return pool.errChan
}

// Alloc try alloc a []byte from internal slab class if no free chunk in slab class Alloc will make one.
func (pool *ChanPool) Alloc(size int) []byte {
	if pool == nil {
		return nil
	}
	if size <= pool.maxSize {
		for i := 0; i < len(pool.classes); i++ {
			if pool.classes[i].size >= size {
				mem := pool.classes[i].pop()
				if mem != nil {
					return mem[:size:size]
				}
				break
			}
		}
	} else {
		len := len(pool.classes)
		if size <= pool.classes[len-1].size {
			mem := pool.classes[len-1].pop()
			if mem != nil {
				return mem[:size:size]
			}
		}
	}
	return make([]byte, size)
}

// Free release a []byte that alloc from Pool.Alloc.
func (pool *ChanPool) Free(mem []byte) {
	if pool == nil {
		return
	}
	size := cap(mem)
	if size <= pool.maxSize {
		for i := 0; i < len(pool.classes); i++ {
			if pool.classes[i].size >= size {
				pool.classes[i].push(mem)
				break
			}
		}
	} else {
		len := len(pool.classes)
		if size <= pool.classes[len-1].size {
			pool.classes[len-1].pop()
		}
	}
	return
}

type chanClass struct {
	size      int
	page      []byte
	pageBegin uintptr
	pageEnd   uintptr
	chunks    chan []byte
	chanPool  *ChanPool
}

func (c *chanClass) push(mem []byte) {
	if c == nil {
		return
	}
	select {
	case c.chunks <- mem:
	default:
		c.chanPool.errChan <- errors.Errorf("size: [%d],  chanClass's channels are overflowing...", c.size)
	}
	return
}

func (c *chanClass) pop() []byte {
	if c == nil {
		return nil
	}
	select {
	case mem := <-c.chunks:
		return mem
	default:
		return nil
	}
}
