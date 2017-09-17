package slab

import "sync"

// SyncPool is a sync.Pool base slab allocation memory pool
type SyncPool struct {
	classes     []sync.Pool
	classesSize []int
	minSize     int
	maxSize     int
}

// newSyncPool create a sync.Pool base slab allocation memory pool.
// minSize is the smallest chunk size.
// maxSize is the lagest chunk size.
// factor is used to control growth of chunk size.
func newSyncPool(minSize, maxSize, factor int) *SyncPool {
	var (
		chunkSize int = minSize
		n         int = 0
	)
	for ; chunkSize <= maxSize; chunkSize *= factor {
		n++
	}
	if maxSize > int(chunkSize/factor) {
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

func (pool *SyncPool) ErrChan() <-chan error {
	return nil
}

// Alloc try alloc a []byte from internal slab class if no free chunk in slab class Alloc will make one.
func (pool *SyncPool) Alloc(size int) []byte {
	if pool == nil {
		return nil
	}
	if size <= pool.classesSize[len(pool.classesSize)-1] {
		for i := 0; i < len(pool.classesSize); i++ {
			if pool.classesSize[i] >= size {
				mem := pool.classes[i].Get().(*[]byte)
				return (*mem)[:size:size]
			}
		}
	}
	return make([]byte, size)
}

// Free release a []byte that alloc from Pool.Alloc.
func (pool *SyncPool) Free(mem []byte) {
	if pool == nil {
		return
	}
	size := cap(mem)
	if size <= pool.classesSize[len(pool.classesSize)-1] {
		for i := 0; i < len(pool.classesSize); i++ {
			if pool.classesSize[i] >= size {
				pool.classes[i].Put(&mem)
				break
			}
		}
	}
	return
}
