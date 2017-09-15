package slab

import (
	"fmt"
	"testing"
	"time"

	"github.com/1046102779/utest"
)

func Test_ChanPool_AllocAndFree(t *testing.T) {
	pool := newChanPool(128, 64*1024, 2)
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

func Test_ChanPool_Alloc_IsNilPtr(t *testing.T) {
	var pool *ChanPool
	mem := pool.Alloc(1024)
	utest.EqualNow(t, cap(mem), 0)
	utest.EqualNow(t, len(mem), 0)
	utest.IsNilNow(t, pool)
}

func Test_ChanPool_Free_IsNilPtr(t *testing.T) {
	var pool *ChanPool
	pool.Free(make([]byte, 64))
	utest.IsNilNow(t, pool)
}

func Test_ChanPool_ErrChan(t *testing.T) {
	pool := newChanPool(128, 1024, 2)
	tick := time.NewTicker(2 * time.Second)
	go func() {
		for {
			select {
			case err := <-pool.ErrChan():
				fmt.Println(err.Error())
				return
			case <-tick.C:
			}
		}
	}()
	return
}
func Test_ChanPool_AllocSmall(t *testing.T) {
	pool := newChanPool(128, 1024, 2)
	mem := pool.Alloc(64)
	utest.EqualNow(t, len(mem), 64)
	utest.EqualNow(t, cap(mem), 128)
	pool.Free(mem)
}

func Test_ChanPool_AllocLarge(t *testing.T) {
	pool := newChanPool(128, 1024, 2)
	mem := pool.Alloc(2048)
	utest.EqualNow(t, len(mem), 2048)
	utest.EqualNow(t, cap(mem), 2048)
	pool.Free(mem)
}

func Test_ChanPool_DoubleFree(t *testing.T) {
	pool := newChanPool(128, 1024, 2)
	mem := pool.Alloc(64)
	go func() {
		defer func() {
			utest.IsNilNow(t, recover())
		}()
		pool.Free(mem)
		pool.Free(mem)
	}()
}

func Test_ChanPool_AllocSlow(t *testing.T) {
	pool := newChanPool(128, 1024, 2)
	mem := pool.classes[len(pool.classes)-1].Pop()
	utest.EqualNow(t, cap(mem), 1024)

	mem = pool.Alloc(1024)
	utest.EqualNow(t, cap(mem), 1024)
}

func Benchmark_ChanPool_AllocAndFree_128(b *testing.B) {
	pool := newChanPool(128, 1024, 2)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			pool.Free(pool.Alloc(128))
		}
	})
}

func Benchmark_ChanPool_AllocAndFree_256(b *testing.B) {
	pool := newChanPool(128, 1024, 2)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			pool.Free(pool.Alloc(256))
		}
	})
}

func Benchmark_ChanPool_AllocAndFree_512(b *testing.B) {
	pool := newChanPool(128, 1024, 2)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			pool.Free(pool.Alloc(512))
		}
	})
}

func Test_ClassPool_AllocAndFree_IsNilPtr(t *testing.T) {
	var class *chanClass
	mem := make([]byte, 64)
	class.Push(mem)
	utest.IsNilNow(t, class)
	tempMem := class.Pop()
	utest.EqualNow(t, cap(tempMem), 0)
	utest.EqualNow(t, cap(tempMem), 0)
}
