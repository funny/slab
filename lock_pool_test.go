package slab

import (
	"testing"

	"github.com/1046102779/utest"
)

func Test_LockPool_ErrChan_NilPtr(t *testing.T) {
	var pool *LockPool
	utest.IsNilNow(t, pool.ErrChan())
}

func Test_LockPool_AllocAndFree(t *testing.T) {
	pool := newLockPool(128, 64*1024, 2)
	for i := 0; i < len(pool.classes); i++ {
		temp := make([][]byte, len(pool.classes[i].chunks))

		for j := 0; j < len(temp); j++ {
			mem := pool.Alloc(pool.classes[i].size)
			utest.EqualNow(t, cap(mem), pool.classes[i].size)
			temp[j] = mem
		}

		for j := 0; j < len(temp); j++ {
			pool.Free(temp[j])
		}
	}
}

func Test_LockPool_AllocSmall(t *testing.T) {
	pool := newLockPool(128, 1024, 2)
	mem := pool.Alloc(64)
	utest.EqualNow(t, len(mem), 64)
	utest.EqualNow(t, cap(mem), 64)
	pool.Free(mem)
}

func Test_LockPool_AllocLarge(t *testing.T) {
	pool := newLockPool(128, 1024, 2)
	mem := pool.Alloc(2048)
	utest.EqualNow(t, len(mem), 2048)
	utest.EqualNow(t, cap(mem), 2048)
	pool.Free(mem)
}

func Test_LockPool_DoubleFree(t *testing.T) {
	pool := newLockPool(128, 1024, 2)
	mem := pool.Alloc(256)
	go func() {
		defer func() {
			utest.NotNilNow(t, recover())
		}()
		pool.Free(mem)
		pool.Free(mem)
	}()
}

func Test_LockPool_AllocSlow(t *testing.T) {
	pool := newLockPool(128, 1024, 2)
	mem := pool.classes[len(pool.classes)-1].pop()
	utest.EqualNow(t, cap(mem), 1024)

	mem = pool.Alloc(1024)
	utest.EqualNow(t, cap(mem), 1024)
}

func Benchmark_LockPool_AllocAndFree_128(b *testing.B) {
	pool := newLockPool(128, 1024, 2)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			pool.Free(pool.Alloc(128))
		}
	})
}

func Benchmark_LockPool_AllocAndFree_256(b *testing.B) {
	pool := newLockPool(128, 1024, 2)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			pool.Free(pool.Alloc(256))
		}
	})
}

func Benchmark_LockPool_AllocAndFree_512(b *testing.B) {
	pool := newLockPool(128, 1024, 2)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			pool.Free(pool.Alloc(512))
		}
	})
}
