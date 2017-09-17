package slab

import (
	"reflect"
	"sync"
	"unsafe"
)

// LockPool is a lock-free slab allocation memory pool.
type LockPool struct {
	classes []lockClass
	minSize int
	maxSize int
}

// newLockPool create a lock-free slab allocation memory pool.
// minSize is the smallest chunk size.
// maxSize is the lagest chunk size.
// factor is used to control growth of chunk size.
// pageSize is the memory size of each slab class.
func newLockPool(minSize, maxSize, factor int, pageSize int) *LockPool {
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
	pool := &LockPool{
		classes: make([]lockClass, n),
		minSize: minSize,
		maxSize: maxSize,
	}

	chunkSize = minSize
	for i := 0; i < n; i++ {
		c := &pool.classes[i]
		c.size = chunkSize
		c.page = make([]byte, pageSize)
		c.chunks = make([][]byte, pageSize/chunkSize)
		c.head = 0
		c.tail = pageSize/chunkSize - 1

		for j := 0; j < len(c.chunks); j++ {
			// lock down the capacity to protect append operation
			c.chunks[j] = c.page[j*chunkSize : (j+1)*chunkSize : (j+1)*chunkSize]
			if j == len(c.chunks)-1 {
				c.pageBegin = uintptr(unsafe.Pointer(&c.page[0]))
				c.pageEnd = uintptr(unsafe.Pointer(&c.chunks[j][0]))
			}
		}

		chunkSize *= factor
	}
	return pool
}

func (pool *LockPool) ErrChan() <-chan error {
	return nil
}

// LockPool try alloc a []byte from internal slab class if no free chunk in slab class Alloc will make one.
func (pool *LockPool) Alloc(size int) []byte {
	if pool == nil {
		return nil
	}
	if size <= pool.classes[len(pool.classes)-1].size {
		for i := 0; i < len(pool.classes); i++ {
			if pool.classes[i].size >= size {
				mem := pool.classes[i].pop()
				if cap(mem) > 0 {
					return mem[:size:size]
				}
			}
		}
	}
	return make([]byte, size)
}

// Free release a []byte that alloc from Pool.Alloc.
func (pool *LockPool) Free(mem []byte) {
	if pool == nil {
		return
	}
	size := cap(mem)
	if size <= pool.classes[len(pool.classes)-1].size {
		for i := 0; i < len(pool.classes); i++ {
			if pool.classes[i].size >= size {
				pool.classes[i].push(mem)
				break
			}
		}
	}
	return
}

type lockClass struct {
	sync.Mutex
	size      int
	page      []byte
	pageBegin uintptr
	pageEnd   uintptr
	chunks    [][]byte
	head      int
	tail      int
}

func (c *lockClass) push(mem []byte) {
	ptr := (*reflect.SliceHeader)(unsafe.Pointer(&mem)).Data
	if c.pageBegin <= ptr && ptr <= c.pageEnd {
		c.Lock()
		c.tail++
		n := c.tail % len(c.chunks)
		// if the panic is received by recover, mutex lock will be deadlock !!
		if c.chunks[n] != nil {
			c.Unlock()
			panic("slab.LockPool: Double Free")
		}
		c.chunks[n] = mem
		c.Unlock()
	}
}

func (c *lockClass) pop() []byte {
	var mem []byte
	c.Lock()
	if c.head <= c.tail {
		n := c.head % len(c.chunks)
		mem = c.chunks[n]
		c.chunks[n] = nil
		c.head++
	}
	c.Unlock()
	return mem
}
