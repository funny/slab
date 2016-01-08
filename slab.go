package slab

import (
	"reflect"
	"sync/atomic"
	"unsafe"
)

// Pool is a lock-free slab allocator.
type Pool struct {
	classes []class
	minSize int
	maxSize int
}

// NewPool create a new memory pool.
// minSize is the smallest chunk size.
// maxSize is the lagest chunk size.
// factor is used to control growth of chunk size.
// pageSize is the memory size of each slab class.
func NewPool(minSize, maxSize, factor, pageSize int) *Pool {
	pool := &Pool{make([]class, 0, 10), minSize, maxSize}
	chunkSize := minSize
	for {
		c := class{
			size:   chunkSize,
			page:   make([]byte, pageSize),
			chunks: make([]chunk, pageSize/chunkSize),
		}
		for i := 0; i < len(c.chunks); i++ {
			chk := &c.chunks[i]
			// lock down the capacity to protect append operation
			chk.mem = c.page[i*chunkSize : (i+1)*chunkSize : (i+1)*chunkSize]
			chk.next = c.head
			c.head = unsafe.Pointer(chk)
		}
		c.beginPtr = uintptr(unsafe.Pointer(&c.chunks[0].mem[0]))
		c.endPtr = uintptr(unsafe.Pointer(&c.chunks[len(c.chunks)-1].mem[0]))
		pool.classes = append(pool.classes, c)

		chunkSize *= factor
		if chunkSize > maxSize {
			break
		}
	}
	return pool
}

// Alloc try alloc a []byte from internal slab class if no free chunk in slab class Alloc will make one.
func (pool *Pool) Alloc(size int) []byte {
	if size <= pool.maxSize {
		capacity := size
		if capacity < pool.minSize {
			capacity = pool.minSize
		}
		for i := 0; i < len(pool.classes); i++ {
			if pool.classes[i].size >= capacity {
				mem := pool.classes[i].Pop()
				if mem != nil {
					return mem[:size]
				}
			}
		}
	}
	return make([]byte, size)
}

// Free release a []byte that alloc from Pool.Alloc.
func (pool *Pool) Free(mem []byte) {
	capacity := cap(mem)
	for i := 0; i < len(pool.classes); i++ {
		if pool.classes[i].size >= capacity {
			pool.classes[i].Push(mem)
		}
	}
}

type class struct {
	size     int
	page     []byte
	chunks   []chunk
	beginPtr uintptr
	endPtr   uintptr
	head     unsafe.Pointer
}

type chunk struct {
	mem  []byte
	next unsafe.Pointer
}

func (c *class) Push(mem []byte) {
	ptr := (*reflect.SliceHeader)(unsafe.Pointer(&mem)).Data
	if c.beginPtr <= ptr && ptr <= c.endPtr {
		chk := &c.chunks[(ptr-c.beginPtr)/uintptr(c.size)]
		if chk.next != nil {
			panic("slab.Pool: Double Free")
		}
		for {
			chk.next = atomic.LoadPointer(&c.head)
			if atomic.CompareAndSwapPointer(&c.head, chk.next, unsafe.Pointer(chk)) {
				break
			}
		}
	}
}

func (c *class) Pop() []byte {
	var ptr unsafe.Pointer
	var chk *chunk
	for {
		ptr = atomic.LoadPointer(&c.head)
		if ptr == nil {
			return nil
		}
		chk = (*chunk)(ptr)
		if atomic.CompareAndSwapPointer(&c.head, ptr, chk.next) {
			break
		}
	}
	chk.next = nil
	return chk.mem
}
