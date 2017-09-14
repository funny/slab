package slab

import "sync"

// SyncPool is a sync.Pool base slab allocation memory pool
type SyncPool struct {
	classes     []sync.Pool
	classesSize []int
	minSize     int
	maxSize     int
}

// NewSyncPool create a sync.Pool base slab allocation memory pool.
// minSize is the smallest chunk size.
// maxSize is the lagest chunk size.
// factor is used to control growth of chunk size.
func NewSyncPool(minSize, maxSize, factor int) *SyncPool {
	var chunkSize int
	n := 0
	for chunkSize = minSize; chunkSize <= maxSize; chunkSize *= factor {
		n++
	}
	if maxSize > chunkSize {
		n++
	}
	pool := &SyncPool{
		make([]sync.Pool, n),
		make([]int, n),
		minSize, maxSize,
	}
	chunkSize = minSize
	for i := 0; i < n; i++ {
		pool.classesSize[i] = chunkSize
		pool.classes[i].New = func(size int) func() interface{} {
			return func() interface{} {
				buf := make([]byte, size)
				return &buf
			}
		}(chunkSize)
		chunkSize *= factor
	}
	return pool
}

// Alloc try alloc a []byte from internal slab class if no free chunk in slab class Alloc will make one.
func (pool *SyncPool) Alloc(size int) []byte {
	if pool == nil {
		return nil
	}
	if size <= pool.maxSize {
		for i := 0; i < len(pool.classesSize); i++ {
			if pool.classesSize[i] >= size {
				mem := pool.classes[i].Get().(*[]byte)
				return (*mem)[:size]
			}
		}
	} else if size > pool.maxSize {
		len := len(pool.classesSize)
		if size <= pool.classesSize[len-1] {
			mem := pool.classes[len-1].Get().(*[]byte)
			return (*mem)[:size]
		}
	}
	return make([]byte, size)
}

// Free release a []byte that alloc from Pool.Alloc.
func (pool *SyncPool) Free(mem []byte) {
	if pool == nil {
		return
	}
	if size := cap(mem); size <= pool.maxSize {
		for i := 0; i < len(pool.classesSize); i++ {
			if pool.classesSize[i] >= size {
				pool.classes[i].Put(&mem)
				return
			}
		}
	} else if size > pool.maxSize {
		len := len(pool.classesSize)
		if size <= pool.classesSize[len-1] {
			pool.classes[len-1].Put(&mem)
			return
		}
	}
	return
}
