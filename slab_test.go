package slab

import (
	"testing"

	"github.com/funny/utest"
)

func Test_AllocFree(t *testing.T) {
	pool := NewPool(128, 64*1024, 2, 1024*1024)
	for i := 0; i < len(pool.classes); i++ {
		temp := make([][]byte, len(pool.classes[i].chunks))

		for j := 0; j < len(temp); j++ {
			mem := pool.Alloc(0, pool.classes[i].size)
			utest.EqualNow(t, cap(mem), pool.classes[i].size)
			temp[j] = mem
		}
		utest.Assert(t, pool.classes[i].head == nil)

		for j := 0; j < len(temp); j++ {
			pool.Free(temp[j])
		}
		utest.NotNilNow(t, pool.classes[i].head)
	}
}

func Test_AllocSmall(t *testing.T) {
	pool := NewPool(128, 1024, 2, 1024)
	mem := pool.Alloc(0, 64)
	utest.EqualNow(t, cap(mem), 128)
}

func Test_AllocLarge(t *testing.T) {
	pool := NewPool(128, 1024, 2, 1024)
	mem := pool.Alloc(0, 2048)
	utest.EqualNow(t, cap(mem), 2048)
}

func Test_DoubleFree(t *testing.T) {
	pool := NewPool(128, 1024, 2, 1024)
	mem := pool.Alloc(0, 64)
	go func() {
		defer func() {
			utest.NotNilNow(t, recover())
		}()
		pool.Free(mem)
		pool.Free(mem)
	}()
}

func Test_AllocSlow(t *testing.T) {
	pool := NewPool(128, 1024, 2, 1024)
	mem := pool.classes[len(pool.classes)-1].Pop()
	utest.EqualNow(t, cap(mem), 1024)
	utest.Assert(t, pool.classes[len(pool.classes)-1].head == nil)

	mem = pool.Alloc(0, 1024)
	utest.EqualNow(t, cap(mem), 1024)
}

func Benchmark_Alloc128_And_Free(b *testing.B) {
	pool := NewPool(128, 1024, 2, 1024)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pool.Free(pool.Alloc(0, 128))
	}
}

func Benchmark_Alloc256_And_Free(b *testing.B) {
	pool := NewPool(128, 1024, 2, 1024)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pool.Free(pool.Alloc(0, 256))
	}
}

func Benchmark_Alloc512_And_Free(b *testing.B) {
	pool := NewPool(128, 1024, 2, 1024)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pool.Free(pool.Alloc(0, 512))
	}
}

func Benchmark_Make128(b *testing.B) {
	var x []byte
	for i := 0; i < b.N; i++ {
		x = make([]byte, 128)
	}
	x = x[:0]
}

func Benchmark_Make256(b *testing.B) {
	var x []byte
	for i := 0; i < b.N; i++ {
		x = make([]byte, 256)
	}
	x = x[:0]
}

func Benchmark_Make512(b *testing.B) {
	var x []byte
	for i := 0; i < b.N; i++ {
		x = make([]byte, 512)
	}
	x = x[:0]
}
