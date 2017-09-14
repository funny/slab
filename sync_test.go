package slab

import (
	"testing"

	"github.com/funny/utest"
)

func Test_SyncPool_NilPtr(t *testing.T) {
	var syncPool *SyncPool
	mem := syncPool.Alloc(64)
	utest.EqualNow(t, cap(mem), 0)
	utest.EqualNow(t, len(mem), 0)
}

func Test_SyncPool_Free_NilPtr(t *testing.T) {
	var syncPool *SyncPool
	mem := make([]byte, 16)
	syncPool.Free(mem)
	// utest project exists a bug for interface{}
	//utest.IsNilNow(t, syncPool)
}

func Test_SyncPool_Alloc_CriticalValue(t *testing.T) {
	pool := NewSyncPool(128, 1000, 2) // test: maxSize > ( last chunkSize = 512 )
	mem := pool.Alloc(1023)
	utest.EqualNow(t, len(mem), 1023)
	utest.EqualNow(t, cap(mem), 1023)
	pool.Free(mem)
}

func Test_SyncPool_AllocSmall_NonIntFactor(t *testing.T) {
	pool := NewSyncPool(128, 1500, 2)
	mem := pool.Alloc(1800)
	utest.EqualNow(t, len(mem), 1800)
	utest.EqualNow(t, cap(mem), 1800)
	pool.Free(mem)
}

func Test_SyncPool_Alloc_NilPtr(t *testing.T) {
	var pool *SyncPool
	mem := pool.Alloc(64)
	utest.EqualNow(t, len(mem), 0)
	utest.EqualNow(t, cap(mem), 0)
}

func Test_SyncPool_AllocSmall(t *testing.T) {
	pool := NewSyncPool(128, 1024, 2)
	mem := pool.Alloc(64)
	utest.EqualNow(t, len(mem), 64)
	utest.EqualNow(t, cap(mem), 128)
	pool.Free(mem)
}

func Test_SyncPool_AllocLarge(t *testing.T) {
	pool := NewSyncPool(128, 1024, 2)
	mem := pool.Alloc(2048)
	utest.EqualNow(t, len(mem), 2048)
	utest.EqualNow(t, cap(mem), 2048)
	pool.Free(mem)
}

func Test_SyncPool_Alloc_BeyondSize(t *testing.T) {
	pool := NewSyncPool(128, 1500, 2)
	mem := pool.Alloc(2500)
	utest.EqualNow(t, len(mem), 2500)
	utest.EqualNow(t, cap(mem), 2500)
	pool.Free(mem)
}
func Test_SyncPool_Alloc_LastElemSize(t *testing.T) {
	pool := NewSyncPool(128, 1500, 2) // maxSize< 1600  <last elem size
	mem := pool.Alloc(1600)
	utest.EqualNow(t, len(mem), 1600)
	utest.EqualNow(t, cap(mem), 1600)
	pool.Free(mem)
}

func Test_SyncPool_AllocLastElem_NonIntFactor(t *testing.T) {
	pool := NewSyncPool(128, 1500, 2)
	mem := pool.Alloc(1800)
	utest.EqualNow(t, len(mem), 1800)
	utest.EqualNow(t, cap(mem), 1800)
	pool.Free(mem)
}

func Benchmark_SyncPool_AllocAndFree_128(b *testing.B) {
	pool := NewSyncPool(128, 1024, 2)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			pool.Free(pool.Alloc(128))
		}
	})
}

func Benchmark_SyncPool_AllocAndFree_256(b *testing.B) {
	pool := NewSyncPool(128, 1024, 2)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			pool.Free(pool.Alloc(256))
		}
	})
}

func Benchmark_SyncPool_AllocAndFree_512(b *testing.B) {
	pool := NewSyncPool(128, 1024, 2)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			pool.Free(pool.Alloc(512))
		}
	})
}

func Benchmark_SyncPool_AllocAndFree_NonIntFactor(b *testing.B) {
	pool := NewSyncPool(128, 1500, 2)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			pool.Free(pool.Alloc(1400))
		}
	})
}

func Benchmark_SyncPool_CacheMiss_128(b *testing.B) {
	pool := NewSyncPool(128, 1024, 2)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			pool.Alloc(128)
		}
	})
}

func Benchmark_SyncPool_CacheMiss_256(b *testing.B) {
	pool := NewSyncPool(128, 1024, 2)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			pool.Alloc(256)
		}
	})
}

func Benchmark_SyncPool_CacheMiss_512(b *testing.B) {
	pool := NewSyncPool(128, 1024, 2)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			pool.Alloc(512)
		}
	})
}

func Benchmark_Make_128(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		var x []byte
		for pb.Next() {
			x = make([]byte, 128)
		}
		x = x[:0]
	})
}

func Benchmark_Make_256(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		var x []byte
		for pb.Next() {
			x = make([]byte, 256)
		}
		x = x[:0]
	})
}

func Benchmark_Make_512(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		var x []byte
		for pb.Next() {
			x = make([]byte, 512)
		}
		x = x[:0]
	})
}
